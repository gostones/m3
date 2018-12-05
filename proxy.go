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
			hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
			b := nb.IsLocal(hostport[0])
			log.Printf("@@@@@ isLocal: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isLocal()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-IPFS-Proxy", "Mirr")

		hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
		port = nb.config.WebPort
		if len(hostport) > 1 {
			port = ParseInt(hostport[1], port)
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
			hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
			b := nb.IsPeer(hostport[0])
			log.Printf("@@@@@ isPeer: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isPeer()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		// resolve peer id
		hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
		id := nb.ToPeerID(hostport[0])
		proxy := nb.GetPeerProxy(id)
		if proxy == "" {
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusServiceUnavailable,
				"No proxy to peer: "+id)
		}

		proxyURL := fmt.Sprintf("http://%v", proxy)
		uri, _ := url.Parse(req.URL.String())
		uri.Host = ToPeerAddr(id)
		if len(hostport) > 1 {
			uri.Host = fmt.Sprintf("%v:%v", uri.Host, hostport[1])
		}

		// copy request
		proxyReq, err := http.NewRequest(req.Method, uri.String(), req.Body)
		if err != nil {
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusInternalServerError,
				"Failed to clone request")
		}

		proxyReq.Header = req.Header
		// for header, values := range req.Header {
		// 	for _, value := range values {
		// 		proxyReq.Header.Add(header, value)
		// 	}
		// }

		proxyReq.Header.Set("Host", req.Host)
		proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

		tr := &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(proxyURL)
		}}

		client := &http.Client{Transport: tr, Timeout: time.Second * 10}

		resp, err := client.Do(proxyReq)

		log.Printf("@@@@@ curl -kv -x %v %v err: %v\n", proxyURL, uri, err)
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
