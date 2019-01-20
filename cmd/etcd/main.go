package main

import (
	"flag"
	"github.com/dhnt/m3/internal"
	"os"
)

var logger = internal.Stdlog

func main() {
	base := flag.String("base", "", "dhnt base")
	flag.Parse()
	if *base == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	es := internal.NewEtcd(*base)
	defer es.Stop()
	es.Run()

	logger.Println("Etcd exited")
}
