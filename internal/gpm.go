package internal

import (
	"context"
	"path/filepath"

	"fmt"
	"github.com/dhnt/m3/internal/misc"
	"github.com/gostones/gpm"

	"io/ioutil"
	"os"
	"os/signal"
)

// ipfs, gogs, mirr
var gpmConfigJSON = `
[
	{
		"name": "ipfs",
		"command": "ipfs daemon",
		"autoRestart": true,
		"workDir": "%v/go/src/github.com/ipfs/go-ipfs"
	},
	{
		"name": "gogs",
		"command": "gogs web --port 3000",
		"autoRestart": true,
		"workDir": "%v/go/src/github.com/gogs/gogs"
	},
	{
		"name": "gotty",
		"command": "gotty --port 50022 --permit-write login",
		"autoRestart": true,
		"workDir": "%v/go/src/github.com/yudai/gotty"
	  },
	{
		"name": "traefik",
		"command": "traefik -c %v/etc/traefik/config.toml --file.directory=%v/etc/traefik",
		"autoRestart": true,
		"workDir": "%v/go/src/github.com/containous/traefik"
	},
	{
		"name": "mirr",
		"command": "mirr --port 18080",
		"autoRestart": true,
		"workDir": "%v/m3"
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

	data = []byte(fmt.Sprintf(gpmConfigJSON, base, base, base, base, base, base, base))
	if err := ioutil.WriteFile(cf, data, 0644); err != nil {
		return "", err
	}
	return string(data), nil
}

type GPM struct {
	signalChan chan bool
}

//
func NewGPM() *GPM {
	return &GPM{
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
	DumpEnv()

	logger.Println("running gpm")

	base := GetDefaultBase()
	if base == "" {
		logger.Println("No DHNT base found!")
		return
	}
	logger.Println("DHNT base:", base)

	// ensure base exist
	if _, err := misc.CreateDir(base); err != nil {
		logger.Println(err)
		return
	}
	//
	pm := gpm.NewProcessManager()
	data, err := readOrCreateConf(base)
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
	signal.Notify(signalChan, os.Interrupt)

	select {
	case err = <-done:
	case <-signalChan:
	case <-r.signalChan:
	}

	logger.Println("Processes terminated")
}

// StartGPM runs gpm server
func StartGPM() {

	s := NewGPM()

	defer s.Stop()

	logger.Printf("starting: %v\n", s)

	s.Run()

	logger.Printf("exited: %v\n", s)
}

// // StartGPM starts core services: p2p, git, and proxy
// func StartGPM() {
// 	base := GetDefaultBase()
// 	if base == "" {
// 		panic("No DHNT base found!")
// 	}
// 	// ensure base exist
// 	if _, err := createDir(base); err != nil {
// 		panic(err)
// 	}
// 	//
// 	pm := gpm.NewProcessManager()
// 	data, err := readOrCreateConf(base)
// 	if err != nil {
// 		panic(err)
// 	}
// 	logger.Println("gpm config: " + data)

// 	err = pm.ParseConfig(data)
// 	if err != nil {
// 		panic(err)
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	done := make(chan error, 1)
// 	go func() {
// 		done <- pm.StartProcesses(ctx)
// 	}()

// 	signalChan := make(chan os.Signal, 1)
// 	signal.Notify(signalChan, os.Interrupt)

// 	select {
// 	case err = <-done:
// 		cancel()
// 		if err != nil {
// 			logger.Println("Error while running processes: ", err)
// 		} else {
// 			logger.Println("Processes finished by themselves.")
// 		}
// 	case <-signalChan:
// 		logger.Println("Got interrupt, stopping processes.")
// 		cancel()
// 		select {
// 		case err = <-done:
// 			if err != nil {
// 				logger.Println("Error while stopping processes: ", err)
// 			} else {
// 				logger.Println("All processes stopped without issues.")
// 			}
// 		}
// 	}
// }
