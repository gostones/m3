package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/kardianos/service"

	m3 "github.com/dhnt/m3/internal"
)

//
var runnerConfigJSON = `
{
	"Name": "m3",
	"DisplayName": "M3",
	"Description": "M3 service",
	
	"Dir": "",
	"Exec": "pm",
	"Args": [],
	"Env": [
		"PATH=%v/go/bin:/bin:/usr/bin",
		"DHNT_BASE=%v",
		"M3_PORT=18080",
		"PORT=18082"
	],
	"Stderr": "%v/tmp/m3d_err.log",
	"Stdout": "%v/tmp/m3d_out.log"
}
`

func readOrCreateConf(base string) ([]byte, error) {
	cf := filepath.Join(base, "etc/m3d.json")
	log.Println("m3 daemon config file: ", cf)

	data, err := ioutil.ReadFile(cf)
	if err == nil {
		return data, nil
	}

	data = []byte(fmt.Sprintf(runnerConfigJSON, base, base, base, base))
	if err := ioutil.WriteFile(cf, data, 0644); err != nil {
		return nil, err
	}
	return data, nil
}

func parseConfig(data []byte) (*Config, error) {
	conf := Config{}
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func getConfig(base string) (*Config, error) {
	data, err := readOrCreateConf(base)
	if err != nil {
		return nil, err
	}
	return parseConfig(data)
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	base := m3.GetDefaultBase()
	if base == "" {
		panic("No DHNT base found!")
	}
	// configPath, err := getConfigPath()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	config, err := getConfig(base)
	log.Println(config)

	if err != nil {
		log.Fatal(err)
	}

	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}

	prg := &program{
		exit: make(chan struct{}),

		Config: config,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
