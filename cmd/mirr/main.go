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

	// var debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	//
	var cfg = &internal.Config{}

	cfg.Port = *port
	cfg.RouteFile = *route

	internal.StartProxy(cfg)
}
