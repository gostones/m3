package main

import (
	"flag"
	"github.com/dhnt/m3/internal"
	"github.com/dhnt/m3/internal/pm"
	"os"
)

var logger = internal.Stdlog

func main() {
	base := flag.String("base", "", "dhnt base")
	port := flag.Int("port", internal.GetDaemonPort(), "gpm port")
	flag.Parse()
	if *base == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	s := pm.NewServer(*base, "", *port)

	defer s.Stop()
	s.Run()

	logger.Println("gpm exited")
}
