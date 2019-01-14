package internal

import (
	"context"
	"path/filepath"

	"fmt"
	"github.com/gostones/gpm"
	"io/ioutil"
	"log"
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
		"name": "mirr",
		"command": "mirr --port 18080",
		"autoRestart": true,
		"workDir": "%v/m3"
	}
]
`

// createDir returns true if base does not exist and was created successfully
// or false if it already exists; otherwise error
func createDir(base string) (bool, error) {
	src, err := os.Stat(base)

	if os.IsNotExist(err) {
		if errDir := os.MkdirAll(base, 0755); errDir != nil {
			return false, errDir
		}
		return true, nil
	}

	if src.Mode().IsRegular() {
		return false, fmt.Errorf("%v exists as file", base)
	}

	return false, nil
}

func readOrCreateConf(base string) (string, error) {
	cf := filepath.Join(base, "etc/gpm.json")
	log.Println("GPM config file: ", cf)

	data, err := ioutil.ReadFile(cf)
	if err == nil {
		return string(data), nil
	}

	data = []byte(fmt.Sprintf(gpmConfigJSON, base, base, base, base))
	if err := ioutil.WriteFile(cf, data, 0644); err != nil {
		return "", err
	}
	return string(data), nil
}

type GPM struct {
	signalChan chan bool
}

func NewGPM() *GPM {
	return &GPM{}
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
	base := GetDefaultBase()
	if base == "" {
		panic("No DHNT base found!")
	}
	// ensure base exist
	if _, err := createDir(base); err != nil {
		panic(err)
	}
	//
	pm := gpm.NewProcessManager()
	data, err := readOrCreateConf(base)
	if err != nil {
		panic(err)
	}
	log.Println("gpm config: " + data)

	err = pm.ParseConfig(data)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- pm.StartProcesses(ctx)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	r.signalChan = make(chan bool, 1)

	select {
	case err = <-done:
		cancel()
		if err != nil {
			log.Println("Error while running processes: ", err)
		} else {
			log.Println("Processes finished by themselves.")
		}
	case <-signalChan:
		log.Println("Got interrupt, stopping processes.")
		cancel()
		select {
		case err = <-done:
			if err != nil {
				log.Println("Error while stopping processes: ", err)
			} else {
				log.Println("All processes stopped without issues.")
			}
		}
	case <-r.signalChan:
		log.Println("Got interrupt, stopping processes.")
		cancel()
		select {
		case err = <-done:
			if err != nil {
				log.Println("Error while stopping processes: ", err)
			} else {
				log.Println("All processes stopped without issues.")
			}
		}
	}
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
// 	log.Println("gpm config: " + data)

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
// 			log.Println("Error while running processes: ", err)
// 		} else {
// 			log.Println("Processes finished by themselves.")
// 		}
// 	case <-signalChan:
// 		log.Println("Got interrupt, stopping processes.")
// 		cancel()
// 		select {
// 		case err = <-done:
// 			if err != nil {
// 				log.Println("Error while stopping processes: ", err)
// 			} else {
// 				log.Println("All processes stopped without issues.")
// 			}
// 		}
// 	}
// }
