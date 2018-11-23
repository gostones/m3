package main

import (
	"flag"
	"log"
)

// Config is application settings
type Config struct {
	Port      int
	ProxyPort int
	Node      Node
}

var config = Config{}

func main() {
	var port = flag.Int("port", 18080, "The port to listen for connections")

	// var debug = flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	config.Port = *port
	config.ProxyPort = FreePort()

	node, err := p2pID()
	if err != nil {
		panic(err)
	}
	config.Node = node

	log.Printf("Configuration port: %v proxy: %v\n", config.Port, config.ProxyPort)

	//
	go loadbalance(config.Port)

	//
	p2pListen(config.ProxyPort)
	httpproxy(config.ProxyPort)
}
