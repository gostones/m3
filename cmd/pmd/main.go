package main

import (
	"github.com/dhnt/m3/internal/daemon"
)

func main() {
	daemon.Run()
}

// import (
// 	"flag"
// 	"fmt"
// 	"github.com/dhnt/m3/internal"
// 	"github.com/dhnt/m3/internal/pm"
// )

// func main() {
// 	host := flag.String("host", "", "service host")
// 	port := flag.Int("port", internal.GetDaemonPort(), "service port")
// 	flag.Parse()

// 	fmt.Println("starting pm service ...")
// 	pm.StartServer(*host, *port)

// 	fmt.Println("pm service stopped")
// }

// func main() {
// 	port := flag.Int("port", internal.GetDaemonPort(), "service port")
// 	m3port := flag.Int("m3port", internal.GetDefaultPort(), "M3 service port")

// 	flag.Parse()

// 	fmt.Println("starting pm http service ...")
// 	pm.StartHTTPServer(*port, *m3port)

// 	fmt.Println("pm http service stopped")
// }
