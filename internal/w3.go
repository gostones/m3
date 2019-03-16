package internal

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"net/http"
)

//W3Proxy start a proxy to W3
func W3Proxy(pid string, port int) {
	address := fmt.Sprintf(":%v", port)
	proxy := goproxy.NewProxyHttpServer()
	proxy.NonproxyHandler = HealthHandlerFunc(fmt.Sprintf("http://127.0.0.1:%v", port))

	proxy.Verbose = true
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if r != nil {
			r.Header.Add("X-W3-Proxy", pid)
			logger.Printf("@@@ OnResponse status: %v length: %v\n", r.StatusCode, r.ContentLength)
		}
		logger.Printf("@@@ W3Proxy OnResponse response: %v\n", r)
		return r
	})

	logger.Printf("W3 proxy listening: %v\n", address)
	logger.Fatal(http.ListenAndServe(address, proxy))
}
