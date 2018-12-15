package main

import (
	"flag"
	"fmt"
	"net/url"
	//"github.com/gostones/mirr/tunnel"
	"log"
	"strings"
)

// Config is application settings
type Config struct {
	//Port      int
	WebPort   int
	ProxyPort int
	ProxyURL  *url.URL
	Local     bool
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
	var port = flag.Int("port", 18080, "The port for http proxy connection")
	var web = flag.Int("web", 80, "The port for traefik reverse proxy connection to local k8s")
	var proxy = flag.String("proxy", "", "Internet firewall http proxy url")
	var local = flag.Bool("local", false, "Allow localhost access")

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
	if *proxy != "" {
		proxyURL, err := url.Parse(*proxy)
		if err == nil {
			cfg.ProxyURL = proxyURL
		}
	}
	cfg.Local = *local
	cfg.WebPort = *web
	cfg.ProxyPort = *port
	//cfg.TunPort = 8022
	//cfg.Pals = pals
	cfg.Aliases = aliases

	// clean up old p2p connections
	err := p2pCloseAll()
	if err != nil {
		panic(err)
	}
	//
	log.Printf("Configuration: %v\n", cfg)

	nb := NewNeighborhood(cfg)

	//
	log.Printf("proxy/p2p port: %v\n", cfg.ProxyPort)
	p2pListen(cfg.ProxyPort)
	httpproxy(cfg.ProxyPort, nb)

	// log.Printf("tunnel port: %v\n", cfg.TunPort)
	// go tunnel.TunServer(cfg.TunPort, "")

	//log.Printf("web/reverse proxy port: %v\n", cfg.WebPort)
	//rpServer(cfg.WebPort, "")
	// loadbalance(cfg.Port, nb)
}
