package main

import (
	"flag"
	"fmt"
	"github.com/gostones/mirr/tunnel"
	"log"
	"strings"
)

// Config is application settings
type Config struct {
	Port      int
	WebPort   int
	ProxyPort int
	ProxyURL  string
	TunPort   int
	Pals      []string
	Aliases   map[string]string
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
	var proxy = flag.String("proxy", "", "http proxy")

	var peers peerFlags
	flag.Var(&peers, "peer", "Peer friends.")

	// var debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	pals := make([]string, len(peers))
	aliases := make(map[string]string)
	for _, v := range peers {
		pa := strings.SplitN(v, ":", 2)
		switch len(pa) {
		case 1:
			pals = append(pals, pa[0])
		case 2:
			aliases[pa[0]] = pa[1]
			pals = append(pals, pa[1])
		}
	}

	//
	var cfg = &Config{}
	cfg.ProxyURL = *proxy
	cfg.Port = *port
	cfg.WebPort = *web
	cfg.ProxyPort = *port //FreePort()
	cfg.TunPort = 8022
	cfg.Pals = pals
	cfg.Aliases = aliases

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

	log.Printf("tunnel port: %v\n", cfg.TunPort)
	go tunnel.TunServer(cfg.TunPort, "")

	log.Printf("proxy port: %v\n", cfg.ProxyPort)
	httpproxy(cfg.ProxyPort, nb)

	// loadbalance(cfg.Port, nb)
}
