package internal

import (
	"context"
	"path/filepath"

	//"fmt"
	"github.com/dhnt/m3/internal/misc"
	"github.com/gostones/gpm"

	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

// ipfs, gogs, mirr
var gpmConfigJSON = `
[
	{
		"name": "etcd",
		"command": "etcd --base ${DHNT_BASE}",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/etcd"
	},
	{
		"name": "ipfs",
		"command": "gsh ${DHNT_BASE}/etc/ipfs/rc.sh",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/ipfs"
	},
	{
		"name": "gogs",
		"command": "gsh ${DHNT_BASE}/etc/gogs/rc.sh",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/gogs"
	},
	{
		"name": "gotty",
		"command": "gotty --port 50022 --permit-write login",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/gotty"
	  },
	{
		"name": "traefik",
		"command": "traefik -c ${DHNT_BASE}/etc/traefik/config.toml --file.directory=${DHNT_BASE}/etc/traefik",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/traefik"
	},
	{
		"name": "mirr",
		"command": "mirr --port 18080",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/m3"
	}
]
`

func readOrCreateConf(base string) (string, error) {
	cf := filepath.Join(base, "etc/gpm.json")
	if _, err := misc.CreateDir(filepath.Dir(cf)); err != nil {
		panic(err)
	}
	logger.Println("GPM config file: ", cf)

	data, err := ioutil.ReadFile(cf)
	if err == nil {
		return string(data), nil
	}

	mapper := func(placeholder string) string {
		switch placeholder {
		case "DHNT_BASE":
			return base
		}
		return ""
	}

	data = []byte(os.Expand(gpmConfigJSON, mapper))
	if err := ioutil.WriteFile(cf, data, 0644); err != nil {
		return "", err
	}
	return string(data), nil
}

type GPM struct {
	base       string
	signalChan chan bool
}

//
func NewGPM(base string) *GPM {
	return &GPM{
		base:       base,
		signalChan: make(chan bool, 1),
	}
}

// Stop stops core services
func (r *GPM) Stop() {
	r.signalChan <- true
}

// Start starts core services: p2p, git, and proxy
func (r *GPM) Start() {
	go r.Run()
}

// Run starts core services
func (r *GPM) Run() {
	logger.Println("running gpm")

	// ensure base exist
	if _, err := misc.CreateDir(r.base); err != nil {
		logger.Println(err)
		return
	}
	if err := os.Chdir(r.base); err != nil {
		logger.Println(err)
		return
	}
	//
	pm := gpm.NewProcessManager()
	data, err := readOrCreateConf(r.base)
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Println("gpm config: " + data)

	err = pm.ParseConfig(data)
	if err != nil {
		logger.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- pm.StartProcesses(ctx)
	}()

	defer cancel()

	signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case err = <-done:
		logger.Println("error:", err)
	case <-signalChan:
	case <-r.signalChan:
	}

	logger.Println("Processes terminated")
}

// StartGPM runs gpm server
func StartGPM(base string) {

	s := NewGPM(base)

	defer s.Stop()

	logger.Printf("starting: %v\n", s)

	s.Run()

	logger.Printf("exited: %v\n", s)
}
