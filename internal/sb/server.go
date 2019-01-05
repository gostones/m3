package sb

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cflager"
	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry-incubator/switchboard/api"
	"github.com/cloudfoundry-incubator/switchboard/apiaggregator"
	"github.com/cloudfoundry-incubator/switchboard/config"
	"github.com/cloudfoundry-incubator/switchboard/domain"
	apirunner "github.com/cloudfoundry-incubator/switchboard/runner/api"
	apiaggregatorrunner "github.com/cloudfoundry-incubator/switchboard/runner/apiaggregator"
	"github.com/cloudfoundry-incubator/switchboard/runner/bridge"
	"github.com/cloudfoundry-incubator/switchboard/runner/health"
	"github.com/cloudfoundry-incubator/switchboard/runner/monitor"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/sigmon"
)

func newConfig(port int, apiport int, be []string) *config.Config {
	var rc config.Config

	rc.API = config.API{
		Port:           uint(apiport),
		AggregatorPort: 0,
		Username:       "admin",
		Password:       "admin",
		ForceHttps:     false,
	}
	rc.Proxy = config.Proxy{
		Port:                     uint(port),
		HealthcheckTimeoutMillis: 5000,
	}

	for _, v := range be {
		hp := strings.SplitN(v, ":", 2)
		p, _ := strconv.ParseUint(hp[1], 10, 32)
		rc.Proxy.Backends = append(rc.Proxy.Backends, config.Backend{
			Name:           v,
			Host:           hp[0],
			Port:           uint(p),
			StatusPort:     uint(p),
			StatusEndpoint: "/health",
		})
	}

	rc.Logger, _ = cflager.New("M3-sb")
	rc.StaticDir = "sb/static"

	return &rc
}

func SwitchBoard(port int, apiport int, be []string) {
	rootConfig := newConfig(port, apiport, be)

	logger := rootConfig.Logger

	// err = rootConfig.Validate()
	// if err != nil {
	// 	logger.Fatal("Error validating config:", err, lager.Data{"config": rootConfig})
	// }

	if _, err := os.Stat(rootConfig.StaticDir); os.IsNotExist(err) {
		logger.Fatal(fmt.Sprintf("staticDir: %s does not exist", rootConfig.StaticDir), nil)
	}

	backends := domain.NewBackends(rootConfig.Proxy.Backends, logger)

	activeNodeClusterMonitor := monitor.NewClusterMonitor(
		backends,
		rootConfig.Proxy.HealthcheckTimeout(),
		logger,
		true,
	)

	activeNodeBridgeRunner := bridge.NewRunner(rootConfig.Proxy.Port, rootConfig.Proxy.ShutdownDelay(), logger)
	clusterStateManager := api.NewClusterAPI(logger)

	activeNodeClusterMonitor.RegisterBackendSubscriber(activeNodeBridgeRunner.ActiveBackendChan)
	activeNodeClusterMonitor.RegisterBackendSubscriber(clusterStateManager.ActiveBackendChan)

	clusterStateManager.RegisterTrafficEnabledChan(activeNodeBridgeRunner.TrafficEnabledChan)
	go clusterStateManager.ListenForActiveBackend()

	apiHandler := api.NewHandler(clusterStateManager, backends, logger, rootConfig.API, rootConfig.StaticDir)
	aggregatorHandler := apiaggregator.NewHandler(logger, rootConfig.API)

	members := grouper.Members{
		{
			Name:   "active-node-bridge",
			Runner: activeNodeBridgeRunner,
		},
		{
			Name:   "api-aggregator",
			Runner: apiaggregatorrunner.NewRunner(rootConfig.API.AggregatorPort, aggregatorHandler),
		},
		{
			Name:   "api",
			Runner: apirunner.NewRunner(rootConfig.API.Port, apiHandler),
		},
		{
			Name:   "active-node-monitor",
			Runner: monitor.NewRunner(activeNodeClusterMonitor, logger),
		},
	}

	if rootConfig.HealthPort != rootConfig.API.Port {
		members = append(members, grouper.Member{
			Name:   "health",
			Runner: health.NewRunner(rootConfig.HealthPort),
		})
	}

	if rootConfig.Proxy.InactiveMysqlPort != 0 {
		inactiveNodeClusterMonitor := monitor.NewClusterMonitor(
			backends,
			rootConfig.Proxy.HealthcheckTimeout(),
			logger,
			false,
		)

		inactiveNodeBridgeRunner := bridge.NewRunner(rootConfig.Proxy.InactiveMysqlPort, rootConfig.Proxy.ShutdownDelay(), logger)

		inactiveNodeClusterMonitor.RegisterBackendSubscriber(inactiveNodeBridgeRunner.ActiveBackendChan)
		clusterStateManager.RegisterTrafficEnabledChan(inactiveNodeBridgeRunner.TrafficEnabledChan)

		members = append(members,
			grouper.Member{
				Name:   "inactive-node-bridge",
				Runner: inactiveNodeBridgeRunner,
			},
			grouper.Member{
				Name:   "inactive-node-monitor",
				Runner: monitor.NewRunner(inactiveNodeClusterMonitor, logger),
			},
		)
	}

	group := grouper.NewOrdered(os.Interrupt, members)
	process := ifrit.Invoke(sigmon.New(group))

	logger.Info("Proxy started", lager.Data{"proxyConfig": rootConfig.Proxy})

	err := <-process.Wait()
	if err != nil {
		logger.Fatal("Switchboard exited unexpectedly", err, lager.Data{"proxyConfig": rootConfig.Proxy})
	}
}
