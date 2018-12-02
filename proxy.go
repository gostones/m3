// https://github.com/elazarl/goproxy
package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"strings"
)

func httpproxy(port int, nb *Neighborhood) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	//
	var isPeer = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			log.Printf("@@@@@ isPeer host: %v\n", req.URL.Host)
			hostPort := strings.Split(req.URL.Host, ":")

			return IsPeerID(hostPort[0])
		}
	}

	proxy.OnRequest(isPeer()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-IPFS-Proxy", "Mirr")

		hostPort := strings.Split(req.URL.Host, ":")
		id := PeerIDB58(hostPort[0])
		b := nb.IsLocal(id)
		host := nb.GetPeerHost(id)

		log.Printf("@@@@@ id: %v local: %v host: %v\n", id, b, host)
		if host == "" {
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusServiceUnavailable,
				"Cannot reach peer: "+id)
		}
		//
		req.URL.Host = host
		req.URL.Scheme = "http"
		req.URL.Host = host
		log.Printf("@@@@@ request modified: %v\n", req)

		return req, nil
	})

	//
	// proxy.OnRequest(goproxy.IsLocalHost).DoFunc(
	// 	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 		return r, goproxy.NewResponse(r,
	// 			goproxy.ContentTypeText, http.StatusForbidden,
	// 			"Don't waste your time!")
	// 	})

	hostport := fmt.Sprintf(":%v", port)
	log.Println("Proxy listening on: " + hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
