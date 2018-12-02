package main

import (
	"flag"
	"log"
)

// Config is application settings
type Config struct {
	Port      int
	WebPort   int
	ProxyPort int
}

func main() {
	var port = flag.Int("port", 18080, "The port to listen for connections")
	var web = flag.Int("web", 8080, "The port to listen for www connections")

	// var debug = flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	var cfg = &Config{}
	cfg.Port = *port
	cfg.WebPort = *web
	cfg.ProxyPort = FreePort()

	log.Printf("Configuration port: %v proxy: %v\n", cfg.Port, cfg.ProxyPort)

	nb := NewNeighborhood(cfg)

	//
	log.Printf("p2p port: %v\n", cfg.ProxyPort)
	p2pListen(cfg.ProxyPort)

	log.Printf("proxy port: %v\n", cfg.ProxyPort)
	go httpproxy(cfg.ProxyPort, nb)

	loadbalance(cfg.Port, nb)
}
