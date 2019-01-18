package internal

import (
	"github.com/coreos/etcd/embed"
	"github.com/dhnt/m3/internal/misc"
	"io/ioutil"
	"path/filepath"
	"time"
)

// EtcdServe starts etcd service
func EtcdServe(base string) {
	file := filepath.Join(base, "etc/etcd.conf.yml")
	dir := filepath.Join(base, "var/default.etcd")
	logger.Println("etcd config file: ", file)

	if err := checkOrCreateEtcdConf(file); err != nil {
		logger.Fatal(err)
	}
	if _, err := misc.CreateDir(dir); err != nil {
		logger.Fatal(err)
	}
	cfg, err := embed.ConfigFromFile(file)
	if err != nil {
		logger.Fatal(err)
	}
	cfg.Dir = dir

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		logger.Printf("Etcd server is ready!")
	case <-time.After(180 * time.Second):
		e.Server.Stop()
		logger.Printf("Etcd server took too long to start!")
	}
	logger.Fatal(<-e.Err())
}

func checkOrCreateEtcdConf(cf string) error {
	if _, err := misc.CreateDir(filepath.Dir(cf)); err != nil {
		return err
	}
	data := []byte(etcdConfigYaml)
	if err := ioutil.WriteFile(cf, data, 0644); err != nil {
		return err
	}
	return nil
}

var etcdConfigYaml = `
# This is the configuration file for the etcd server.

# Human-readable name for this member.
name: 'default'

# Path to the data directory.
data-dir:

# Path to the dedicated wal directory.
wal-dir:

# Number of committed transactions to trigger a snapshot to disk.
snapshot-count: 10000

# Time (in milliseconds) of a heartbeat interval.
heartbeat-interval: 100

# Time (in milliseconds) for an election to timeout.
election-timeout: 1000

# Raise alarms when backend size exceeds the given quota. 0 means use the
# default quota.
quota-backend-bytes: 0

# List of comma separated URLs to listen on for peer traffic.
listen-peer-urls: http://localhost:2380

# List of comma separated URLs to listen on for client traffic.
listen-client-urls: http://localhost:2379

# Maximum number of snapshot files to retain (0 is unlimited).
max-snapshots: 5

# Maximum number of wal files to retain (0 is unlimited).
max-wals: 5

# Comma-separated white list of origins for CORS (cross-origin resource sharing).
cors:

# List of this member's peer URLs to advertise to the rest of the cluster.
# The URLs needed to be a comma-separated list.
initial-advertise-peer-urls: http://localhost:2380

# List of this member's client URLs to advertise to the public.
# The URLs needed to be a comma-separated list.
advertise-client-urls: http://localhost:2379

# Discovery URL used to bootstrap the cluster.
discovery:

# Valid values include 'exit', 'proxy'
discovery-fallback: 'proxy'

# HTTP proxy to use for traffic to discovery service.
discovery-proxy:

# DNS domain used to bootstrap initial cluster.
discovery-srv:

# Initial cluster configuration for bootstrapping.
initial-cluster:

# Initial cluster token for the etcd cluster during bootstrap.
initial-cluster-token: 'etcd-cluster'

# Initial cluster state ('new' or 'existing').
initial-cluster-state: 'new'

# Reject reconfiguration requests that would cause quorum loss.
strict-reconfig-check: false

# Accept etcd V2 client requests
enable-v2: true

# Enable runtime profiling data via HTTP server
enable-pprof: true

# Valid values include 'on', 'readonly', 'off'
proxy: 'off'

# Time (in milliseconds) an endpoint will be held in a failed state.
proxy-failure-wait: 5000

# Time (in milliseconds) of the endpoints refresh interval.
proxy-refresh-interval: 30000

# Time (in milliseconds) for a dial to timeout.
proxy-dial-timeout: 1000

# Time (in milliseconds) for a write to timeout.
proxy-write-timeout: 5000

# Time (in milliseconds) for a read to timeout.
proxy-read-timeout: 0

client-transport-security:
  # Path to the client server TLS cert file.
  cert-file:

  # Path to the client server TLS key file.
  key-file:

  # Enable client cert authentication.
  client-cert-auth: false

  # Path to the client server TLS trusted CA cert file.
  trusted-ca-file:

  # Client TLS using generated certificates
  auto-tls: false

peer-transport-security:
  # Path to the peer server TLS cert file.
  cert-file:

  # Path to the peer server TLS key file.
  key-file:

  # Enable peer client cert authentication.
  client-cert-auth: false

  # Path to the peer server TLS trusted CA cert file.
  trusted-ca-file:

  # Peer TLS using generated certificates.
  auto-tls: false

# Enable debug-level logging for etcd.
debug: false

logger: zap

# Specify 'stdout' or 'stderr' to skip journald logging even when running under systemd.
log-outputs: [stderr]

# Force to create a new one member cluster.
force-new-cluster: false

auto-compaction-mode: periodic
auto-compaction-retention: "1"
`
