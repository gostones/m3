package main

import (
	"flag"
	"strings"

	internal "github.com/dhnt/m3/internal"
)

func main() {
	var port = flag.Int("port", 18080, "Bind port")

	var home internal.ListFlags
	flag.Var(&home, "home", "Home routes domain/host:port")

	var proxy internal.ListFlags
	flag.Var(&proxy, "proxy", "Peer ID as Internet proxy")

	var local = flag.Bool("local", false, "Allow localhost access")
	var blocked internal.ListFlags
	flag.Var(&blocked, "block", "Block port if local is enabled")

	//
	var alias internal.ListFlags
	flag.Var(&alias, "alias", "Peer ID alias name:id")

	// var debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	aliasMap := make(map[string]string)
	for _, v := range alias {
		pa := strings.SplitN(v, ":", 2)
		switch len(pa) {
		case 1:
			// no alias
		case 2:
			aliasMap[pa[0]] = pa[1]
		}
	}

	homereg := internal.NewRouteRegistry()
	homereg.SetDefault("localhost:80")
	for _, v := range home {
		pa := strings.SplitN(v, "/", 2) // domain/host:port
		switch len(pa) {
		case 1:
			// invalid
		case 2:
			if pa[0] == "" {
				homereg.SetDefault(pa[1])
			} else {
				homereg.Add(pa[0], pa[1])
			}
		}
	}

	//
	var cfg = &internal.Config{}

	cfg.Port = *port
	cfg.Local = *local
	cfg.Home = homereg
	cfg.Blocked = blocked
	cfg.Proxy = proxy
	cfg.Alias = aliasMap

	internal.StartProxy(cfg)
}
