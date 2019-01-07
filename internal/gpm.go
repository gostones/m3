package internal

import (
	"context"
	"fmt"
	"github.com/gostones/gpm"
	"log"
	"os"
	"os/signal"
)

var gpm_config_json = `
[
  {
    "name": "gogs",
    "command": "gogs web --port %v"
  },
  {
    "name": "ipfs",
    "command": "ipfs daemon",
    "autoRestart": true
  }
]
`

func StartGPM(gitPort int) {
	pm := gpm.NewProcessManager()
	data := fmt.Sprintf(gpm_config_json, gitPort)
	log.Println("gpm config: " + data)
	err := pm.ParseConfig(data)
	if err != nil {
		log.Println("Could not parse config file", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- pm.StartProcesses(ctx)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

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
	}
}
