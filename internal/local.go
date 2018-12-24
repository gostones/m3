package internal

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

//LocalProxy start a local proxy to w3
func LocalProxy(port int) {
	hostport := fmt.Sprintf("localhost:%v", port)
	proxy := goproxy.NewProxyHttpServer()
	//proxy.ConnectDial = nil
	proxy.Verbose = true
	log.Printf("local proxy listening on: %v\n", hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
