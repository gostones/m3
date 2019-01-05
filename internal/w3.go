package internal

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

//W3Proxy start a proxy to W3
func W3Proxy(port int) {
	hostport := fmt.Sprintf(":%v", port)
	proxy := goproxy.NewProxyHttpServer()
	proxy.NonproxyHandler = HealthHandlerFunc(fmt.Sprintf("http://127.0.0.1:%v", port))

	proxy.Verbose = true

	log.Printf("W3 proxy listening: %v\n", hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
