// https://github.com/elazarl/goproxy
package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func httpproxy(port int, nb *Neighborhood) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	//
	var isLocal = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			hostPort := strings.Split(req.URL.Host, ":")
			id := PeerIDB58(hostPort[0])
			b := id != "" && nb.IsLocal(id)
			log.Printf("@@@@@ isLocal: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isLocal()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-IPFS-Proxy", "Mirr")

		hostPort := strings.Split(req.URL.Host, ":")
		port = nb.config.WebPort
		if len(hostPort) > 1 {
			port = ParseInt(hostPort[1], port)
		}
		host := fmt.Sprintf("localhost:%v", port)

		//
		req.URL.Host = host
		req.URL.Scheme = "http"
		req.URL.Host = host
		log.Printf("@@@@@ local request modified: %v\n", req)

		return req, nil
	})

	//
	var isPeer = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			hostPort := strings.Split(req.URL.Host, ":")
			id := PeerIDB58(hostPort[0])
			b := id != "" && !nb.IsLocal(id)
			log.Printf("@@@@@ isPeer: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isPeer()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		//copy request
		proxyReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
		if err != nil {
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusInternalServerError,
				"failed to clone request")
		}

		proxyReq.Header = req.Header
		// for header, values := range req.Header {
		// 	for _, value := range values {
		// 		proxyReq.Header.Add(header, value)
		// 	}
		// }

		proxyReq.Header.Set("Host", req.Host)
		proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

		//
		hostPort := strings.Split(req.URL.Host, ":")
		id := PeerIDB58(hostPort[0])
		proxyURL := "http://" + nb.GetPeerHost(id)
		tr := &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(proxyURL)
		}}

		client := &http.Client{Transport: tr, Timeout: time.Second * 10}

		resp, err := client.Do(proxyReq)

		log.Printf("@@@@@ curl -kv -x %v %v err: %v\n", proxyURL, req.URL, err)
		if err != nil {
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusServiceUnavailable,
				"Cannot reach peer: "+id)
		}

		return req, resp
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
