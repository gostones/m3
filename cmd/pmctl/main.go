package main

import (
	//"flag"
	"fmt"
	"github.com/dhnt/m3/internal"
	"github.com/dhnt/m3/internal/pm"
	"os"
)

func main() {
	usage := "usage: pmctl start|stop|status"
	// TODO
	// host := flag.String("host", "localhost", "m3d service host")
	// port := flag.Int("port", internal.GetDaemonPort(), "m3d service port")

	// flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println(usage)
		return
	}

	cli, err := pm.NewClient("localhost", internal.GetDaemonPort())
	if err != nil {
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "start":
		cli.Start()
	case "stop":
		cli.Stop()
	case "status":
		cli.Status()
	default:
		fmt.Println(usage)
	}
}
