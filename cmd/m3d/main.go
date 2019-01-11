package main

import (
	"fmt"
	"github.com/dhnt/m3/internal/daemon"
	"os"
)

func main() {
	usage := "Usage: m3d install | remove | start | stop | status"

	if len(os.Args) != 2 {
		fmt.Println(usage)
	}
	daemon.Startup()
}
