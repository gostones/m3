// https://github.com/elazarl/goproxy
package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	//"time"
	//"github.com/gostones/mirr/tunnel"
)

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(u)

	// Update the headers to allow for SSL redirection
	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = u.Host

	proxy.ServeHTTP(res, req)
}

func httpproxy(port int, nb *Neighborhood) {

	log.Printf("@@@@@ ProxyURL: %v\n", nb.config.ProxyURL)

	//
	proxy := goproxy.NewProxyHttpServer()
	proxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		log.Printf("@@@@@ Proxy: %v\n", req.URL)

		hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
		local := IsLocalHost(hostport[0])
		if local {
			log.Printf("@@@@@ Proxy local: %v\n", local)
			return nil, nil
		}
		home := nb.IsLocal(hostport[0])
		if home {
			log.Printf("@@@@@ Proxy home: %v\n", home)
			return nil, nil
		}

		peer := nb.IsPeer(hostport[0])
		if peer {
			log.Printf("@@@@@ Proxy peer: %v\n", peer)

			port := 80
			if len(hostport) > 1 {
				port = ParseInt(hostport[1], port)
			}

			// resolve peer id
			id := nb.ToPeerID(hostport[0])
			//addr := ToPeerAddr(id)

			proxy := nb.GetPeerProxy(id)
			if proxy == "" {
				return nil, nil
			}
			return url.Parse(fmt.Sprintf("http://%v", proxy))
		}

		log.Printf("@@@@@ Proxy default url: %v\n", nb.config.ProxyURL)

		return nb.config.ProxyURL, nil
	}
	proxy.ConnectDial = nil
	if nb.config.ProxyURL != nil {
		proxy.ConnectDial = proxy.NewConnectDialToProxy(nb.config.ProxyURL.String())
	}

	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		log.Printf("@@@@@ Dial network: %v addr: %v\n", network, addr)

		hostport := strings.Split(strings.ToLower(addr), ":")

		//
		if IsLocalHost(hostport[0]) {
			log.Printf("@@@@@ Dial blocking local. network: %v host: %v\n", network, addr)

			c, err = nil, fmt.Errorf("Nice try: %v", hostport)
			return
		}

		if nb.IsLocal(hostport[0]) {
			log.Printf("@@@@@ Dial home. network: %v host: %v\n", network, addr)

			addr = fmt.Sprintf("127.0.0.1:%v", nb.config.WebPort)
			c, err = net.Dial(network, addr)
			return
		}

		//
		peer := nb.IsPeer(hostport[0])
		if peer {
			log.Printf("@@@@@ Dial peer . network: %v addr: %v peer: %v\n", network, addr, peer)

			// resolve peer id
			id := nb.ToPeerID(hostport[0])

			proxy := nb.GetPeerProxy(id)
			if proxy == "" {
				c, err = nil, fmt.Errorf("Peer addr not found: %v", hostport)
				return
			}
		}

		log.Printf("@@@@@ Dial web . network: %v addr: %v\n", network, addr)

		//c, err = nil, fmt.Errorf("Internet is blocked: %v", hostport)

		c, err = net.Dial(network, addr)
		return
	}

	proxy.Verbose = true

	// non proxy request handling
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", 500)
	})

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("@@@@@ on request: %v\n", req)
			return req, nil
		})

	///////////
	// // block localhost
	// var isLocalHost = func() goproxy.ReqConditionFunc {
	// 	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
	// 		hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
	// 		b := IsLocalHost(hostport[0])
	// 		log.Printf("@@@@@ isLocalHost: %v host: %v\n", b, req.URL.Host)
	// 		return b
	// 	}
	// }
	// proxy.OnRequest(isLocalHost()).DoFunc(
	// 	func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 		if nb.config.Local {
	// 			return req, nil
	// 		}
	// 		return req, goproxy.NewResponse(req,
	// 			goproxy.ContentTypeText, http.StatusForbidden,
	// 			"Don't waste your time!")
	// 	})

	// // local home node - forward to k8s
	// var isHome = func() goproxy.ReqConditionFunc {
	// 	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
	// 		hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
	// 		b := nb.IsLocal(hostport[0])
	// 		log.Printf("@@@@@ isHome: %v host: %v\n", b, req.URL.Host)
	// 		return b
	// 	}
	// }
	// proxy.OnRequest(isHome()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// 	req.Header.Set("X-IPFS-Proxy", "Mirr")

	// 	hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
	// 	target := fmt.Sprintf("127.0.0.1:%v", nb.config.WebPort)

	// 	//
	// 	req.Host = hostport[0] //discard port
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = target
	// 	log.Printf("@@@@@ local request modified: %v\n", req)

	// 	return req, nil
	// })

	// // tunnel to peer
	// var isPeer = func() goproxy.ReqConditionFunc {
	// 	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
	// 		hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
	// 		b := nb.IsPeer(hostport[0])
	// 		log.Printf("@@@@@ isPeer: %v host: %v\n", b, req.URL.Host)
	// 		return b
	// 	}
	// }
	// //var tunnels = make(map[string]int)
	// proxy.OnRequest(isPeer()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// 	req.Header.Set("X-IPFS-Proxy", "Mirr")
	// 	//

	// 	// tunURL := fmt.Sprintf("http://%v:%v", addr, nb.config.TunPort)
	// 	// host := fmt.Sprintf("%v:%v", addr, port)
	// 	// locPort, ok := tunnels[host]
	// 	// if !ok {
	// 	// 	locPort := FreePort()
	// 	// 	remote := fmt.Sprintf("localhost:%v:localhost:%v", locPort, port)
	// 	// 	go tunnel.TunClient(proxyURL, tunURL, remote)
	// 	// 	tunnels[host] = locPort
	// 	// 	log.Printf("@@@@@ peer remote: %v\n", remote)
	// 	// }

	// 	hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
	// 	port := 80
	// 	if len(hostport) > 1 {
	// 		port = ParseInt(hostport[1], port)
	// 	}

	// 	// resolve peer id
	// 	id := nb.ToPeerID(hostport[0])
	// 	addr := ToPeerAddr(id)
	// 	host := fmt.Sprintf("%v:%v", addr, port)

	// 	//
	// 	req.Host = host
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = host
	// 	log.Printf("@@@@@ peer request modified: %v\n", req)

	// 	return req, nil
	// })

	///////////////
	// proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {
	// 		ctx.UserData, resp, err = tr.DetailedRoundTrip(req)
	// 		return
	// 	})
	// 	logger.LogReq(req, ctx)
	// 	return req, nil
	// })

	// proxy.OnRequest(isPeer()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	// resolve peer id
	// 	hostport := strings.Split(strings.ToLower(req.URL.Host), ":")
	// 	id := nb.ToPeerID(hostport[0])
	// 	proxy := nb.GetPeerProxy(id)
	// 	if proxy == "" {
	// 		return req, goproxy.NewResponse(req,
	// 			goproxy.ContentTypeText, http.StatusServiceUnavailable,
	// 			"No proxy to peer: "+id)
	// 	}

	// 	proxyURL := fmt.Sprintf("http://%v", proxy)
	// 	uri, _ := url.Parse(req.URL.String())
	// 	uri.Host = ToPeerAddr(id)
	// 	if len(hostport) > 1 {
	// 		uri.Host = fmt.Sprintf("%v:%v", uri.Host, hostport[1])
	// 	}

	// 	// copy request
	// 	proxyReq, err := http.NewRequest(req.Method, uri.String(), req.Body)
	// 	if err != nil {
	// 		return req, goproxy.NewResponse(req,
	// 			goproxy.ContentTypeText, http.StatusInternalServerError,
	// 			"Failed to clone request")
	// 	}

	// 	proxyReq.Header = req.Header
	// 	// for header, values := range req.Header {
	// 	// 	for _, value := range values {
	// 	// 		proxyReq.Header.Add(header, value)
	// 	// 	}
	// 	// }

	// 	proxyReq.Header.Set("Host", req.Host)
	// 	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	// 	tr := &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
	// 		return url.Parse(proxyURL)
	// 	}}

	// 	client := &http.Client{Transport: tr, Timeout: time.Second * 10}

	// 	resp, err := client.Do(proxyReq)

	// 	log.Printf("@@@@@ curl -kv -x %v %v err: %v\n", proxyURL, uri, err)
	// 	if err != nil {
	// 		return req, goproxy.NewResponse(req,
	// 			goproxy.ContentTypeText, http.StatusServiceUnavailable,
	// 			"Cannot reach peer: "+id)
	// 	}

	// 	return req, resp
	// })

	//proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	// proxy.OnRequest().DoFunc(func (req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 	log.Printf("@@@@@ www request: %v\n", req)

	// 	if req.URL.Scheme == "https" {
	// 		req.URL.Scheme = "http"
	// 	}
	// 	return req, nil
	// })

	hostport := fmt.Sprintf(":%v", port)
	log.Println("Proxy listening on: " + hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
