// https://github.com/elazarl/goproxy
package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func httpproxy(port int) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	//
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-IPFS-Proxy", "Mirr")
			return r, nil
		})
	//
	var isPeer = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			return IsPeerID(req.URL.Host)
		}
	}
	proxy.OnRequest(isPeer()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusForbidden,
				"Can't connect to peer! "+r.URL.Host)
		})

	proxy.OnRequest(goproxy.IsLocalHost).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusForbidden,
				"Don't waste your time!")
		})

	// proxy.OnRequest(goproxy.DstHostIs("www.reddit.com")).DoFunc(
	// 	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 		h, _, _ := time.Now().Clock()
	// 		if h >= 8 && h <= 17 {
	// 			return r, goproxy.NewResponse(r,
	// 				goproxy.ContentTypeText, http.StatusForbidden,
	// 				"Don't waste your time!")
	// 		} else {
	// 			ctx.Warnf("clock: %d, you can waste your time...", h)
	// 		}
	// 		return r, nil
	// 	})

	//
	// proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).
	// 	HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
	// 	defer func() {
	// 		if e := recover(); e != nil {
	// 			ctx.Logf("error connecting to remote: %v", e)
	// 			client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
	// 		}
	// 		client.Close()
	// 	}()
	// 	clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	// 	remote, err := net.Dial("tcp", req.URL.Host)
	// 	orPanic(err)
	// 	remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
	// 	for {
	// 		req, err := http.ReadRequest(clientBuf.Reader)
	// 		orPanic(err)
	// 		orPanic(req.Write(remoteBuf))
	// 		orPanic(remoteBuf.Flush())
	// 		resp, err := http.ReadResponse(remoteBuf.Reader, req)
	// 		orPanic(err)
	// 		orPanic(resp.Write(clientBuf.Writer))
	// 		orPanic(clientBuf.Flush())
	// 	}
	// })
	//

	hostport := fmt.Sprintf(":%v", port)
	log.Println("Proxy listening on: " + hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
