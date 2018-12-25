package main

import (
	"flag"
	"strings"

	internal "github.com/dhnt/m3/internal"
)

func main() {
	var port = flag.Int("port", 18080, "The port for http proxy connection")

	var web = flag.String("web", "localhost:80", "The web host:port for traefik reverse proxy connection to local home k8s")
	//var proxy = flag.String("proxy", "", "Internet firewall http proxy url")
	var local = flag.Bool("local", false, "Allow localhost access")

	//
	var blocked internal.ListFlags
	flag.Var(&blocked, "block", "Silently disregard requests from specified ports")

	//
	var peers internal.ListFlags
	flag.Var(&peers, "peer", "Peer friends.")

	// var debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	pals := []string{} //make([]string, len(peers))
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
	var cfg = &internal.Config{}

	cfg.Port = *port
	cfg.Local = *local
	cfg.WebHost = *web
	cfg.Blocked = blocked
	cfg.Pals = pals
	cfg.Aliases = aliases

	internal.StartProxy(cfg)
}
