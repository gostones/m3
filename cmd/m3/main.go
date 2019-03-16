package main

import (
	"github.com/dhnt/m3/internal"
	"github.com/dhnt/m3/internal/daemon"
)

var logger = internal.Logger()

func main() {
	logger.Info("starting m3 daemon ...")
	daemon.Run()
}
