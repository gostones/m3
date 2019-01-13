package main

import (
	m3 "github.com/dhnt/m3/internal"

	"github.com/dhnt/m3/internal/daemon"
)

func main() {
	m3.SetDefaultEnv()

	daemon.Startup()
}
