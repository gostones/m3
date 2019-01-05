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
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if r != nil {
			r.Header.Add("X-W3-Proxy", hostport)
			log.Printf("@@@ OnResponse status: %v length: %v\n", r.StatusCode, r.ContentLength)
		}
		log.Printf("@@@ W3Proxy OnResponse response: %v\n", r)
		return r
	})

	log.Printf("W3 proxy listening: %v\n", hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
