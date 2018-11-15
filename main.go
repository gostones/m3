package main

import (
	"flag"
)

func main() {
	var port = flag.Int("port", 18080, "The port to listen for connections")
	var pport = flag.Int("proxy", 10080, "The port to listen for connections")

	// var debug = flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	//export http_proxy=http://localhost:18080
	be := []string{} //"localhost:50081", "localhost:50082"}
	go loadbalance(*port, be)

	//
	// go forward("localhost:50081", "localhost:10080")
	// go forward("localhost:50082", "localhost:10080")

	//
	httpproxy(*pport)
}
