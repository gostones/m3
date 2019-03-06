package internal

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/dhnt/m3/internal/misc"
	"github.com/gostones/gpm"
)

//
// port assignment:
// 0xxxx - gost
// 1xxxx - mirr
// 2xxxx - traefik
// 3xxxx - frps/chisel
// 4xxxx - frpc/chisel
// 5xxxx - other
// 6xxxx - free
//
var gpmConfigJSON = `
[
	{
		"name": "etcd",
		"command": "etcd --config-file ${DHNT_BASE}/etc/etcd.conf.yml",
		"autoRestart": true
	},
	{
		"name": "ipfs",
		"command": "gsh ${DHNT_BASE}/etc/ipfs-rc.sh",
		"autoRestart": true
	},
	{
		"name": "gogs",
		"command": "gsh ${DHNT_BASE}/etc/gogs-rc.sh",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}/home/gogs"
	},
	{
		"name": "gotty",
		"command": "gotty --port 50022 --permit-write login",
		"autoRestart": true
	},
	{
		"name": "fb",
		"command": "gsh ${DHNT_BASE}/etc/fb-rc.sh",
		"autoRestart": true
	},
	{
		"name": "caddy",
		"command": "caddy -conf ${DHNT_BASE}/etc/Caddyfile",
		"autoRestart": true
	},
	{
		"name": "frps",
		"command": "frps -c ${DHNT_BASE}/etc/frps.ini",
		"autoRestart": true
	},
	{
		"name": "traefik",
		"command": "traefik -c ${DHNT_BASE}/etc/traefik/traefik.toml",
		"autoRestart": true,
		"workDir": "${DHNT_BASE}"
	},
	{
		"name": "mirr",
		"command": "mirr --port 18080 --route ${DHNT_BASE}/etc/route.conf",
		"autoRestart": true
	},
	{
		"name": "gost",
		"command": "gost -L http://:8080  -L https://:8443 -L socks5://:1080 -F http://127.0.0.1:18080",
		"autoRestart": true
	},
	{
		"name": "chisel",
		"command": "chisel server --port 38080",
		"autoRestart": true
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
	if err != nil {
		//create default
		data = []byte(gpmConfigJSON)
		if err := ioutil.WriteFile(cf, data, 0644); err != nil {
			return "", err
		}
	}

	mapper := func(placeholder string) string {
		switch placeholder {
		case "DHNT_BASE":
			return base
		}
		return ""
	}
	return os.Expand(string(data), mapper), nil
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
