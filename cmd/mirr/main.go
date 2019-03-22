package main

import (
	"flag"

	"github.com/dhnt/m3/internal"
)

var logger = internal.Logger()

func main() {
	//
	var port = flag.Int("port", 18080, "Bind port")
	var route = flag.String("route", "route.conf", "Route configuration")

	// var debug = flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	//
	var cfg = &internal.Config{}

	cfg.Port = *port
	cfg.RouteFile = *route

	logger.Info("starting mirr ...")
	logger.Infof("configration: %v", cfg)

	internal.StartProxy(cfg)
}
