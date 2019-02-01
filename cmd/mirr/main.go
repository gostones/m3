package main

import (
	"flag"
	// "strings"

	"github.com/dhnt/m3/internal"
)

func main() {
	//
	var port = flag.Int("port", internal.GetDefaultPort(), "Bind port")
	var route = flag.String("route", "route.conf", "Route configuration")

	// var home m3.ListFlags
	// flag.Var(&home, "home", "Home routes domain/host:port")

	// var web m3.ListFlags
	// flag.Var(&web, "web", "Peer ID as W3 proxy ")

	// var local = flag.Bool("local", true, "Allow localhost access")
	// var blocked m3.ListFlags
	// flag.Var(&blocked, "block", "Block port if local is enabled")

	// //
	// var alias m3.ListFlags
	// flag.Var(&alias, "alias", "Peer ID alias name:id")

	// var debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// aliasMap := make(map[string]string)
	// for _, v := range alias {
	// 	pa := strings.SplitN(v, ":", 2)
	// 	switch len(pa) {
	// 	case 1:
	// 		// no alias
	// 	case 2:
	// 		aliasMap[pa[0]] = pa[1]
	// 	}
	// }

	//
	var cfg = &internal.Config{}

	cfg.Port = *port
	cfg.RouteFile = *route

	// cfg.Local = *local
	// cfg.Blocked = blocked
	// cfg.Home = home
	// cfg.Web = web
	// cfg.Alias = aliasMap

	internal.StartProxy(cfg)
}
