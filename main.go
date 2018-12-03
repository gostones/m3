package main

import (
	"flag"
	"fmt"
	"log"
)

// Config is application settings
type Config struct {
	Port      int
	WebPort   int
	ProxyPort int
	Pals      []string
}

type peerFlags []string

func (i *peerFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *peerFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var port = flag.Int("port", 18080, "The port to listen for connections")
	var web = flag.Int("web", 8080, "The port to listen for www connections")
	var pals peerFlags
	flag.Var(&pals, "peer", "Peer friends.")

	// var debug = flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	var cfg = &Config{}
	cfg.Port = *port
	cfg.WebPort = *web
	cfg.ProxyPort = FreePort()
	cfg.Pals = pals

	//
	err := p2pCloseAll()
	if err != nil {
		panic(err)
	}
	//
	log.Printf("Configuration: %v\n", cfg)

	nb := NewNeighborhood(cfg)

	//
	log.Printf("p2p port: %v\n", cfg.ProxyPort)
	p2pListen(cfg.ProxyPort)

	log.Printf("proxy port: %v\n", cfg.ProxyPort)
	go httpproxy(cfg.ProxyPort, nb)

	loadbalance(cfg.Port, nb)
}
