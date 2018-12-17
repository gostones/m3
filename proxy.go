// https://github.com/elazarl/goproxy
package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	//"net/http/httputil"
	//"github.com/vulcand/oxy/forward"
	"crypto/tls"
	"net/url"
	"strings"
	//"time"
	//"github.com/gostones/mirr/tunnel"
)

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	//u, _ := url.Parse(target)
	//req.URL = testutils.ParseURI("http://localhost:63450")
	//fwd.ServeHTTP(w, req)

	// proxy := httputil.NewSingleHostReverseProxy(u)

	// // Update the headers to allow for SSL redirection
	// req.URL.Host = u.Host
	// req.URL.Scheme = u.Scheme
	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// req.Host = u.Host

	// proxy.ServeHTTP(res, req)
}

func httpproxy(port int, nb *Neighborhood) {

	log.Printf("@@@@@ ProxyURL: %v\n", nb.config.ProxyURL)

	//
	proxy := goproxy.NewProxyHttpServer()

	p := func(req *http.Request) (*url.URL, error) {
		log.Printf("@@@@@ Proxying: %v %v %v url: %v\n", req.Proto, req.Method, req.Host, req.URL)

		hostport := strings.Split(req.URL.Host, ":")
		proxyURL := nb.config.ProxyURL

		if IsLocalHost(hostport[0]) || IsHome(hostport[0]) {
			return nil, nil
		}

		if IsPeer(hostport[0]) {
			log.Printf("@@@@@ Proxy url: %v\n", req.URL)

			tld := TLD(hostport[0])
			id := ToPeerID(tld)
			if id == "" {
				return nil, fmt.Errorf("Peer invalid: %v", hostport[0])
			}
			target := nb.GetPeerProxy(id)
			if target == "" {
				return nil, fmt.Errorf("Peer not reachable: %v", hostport[0])
			}

			log.Printf("@@@@@ Proxy peer url: %v target: %v\n", req.URL, target)

			proxyURL, _ = url.Parse(fmt.Sprintf("http://%v", target))
			return proxyURL, nil
		}

		return proxyURL, nil
	}

	dial := func(network, addr string) (net.Conn, error) {
		hostport := strings.Split(addr, ":")

		if IsHome(hostport[0]) {
			target := fmt.Sprintf("127.0.0.1:%v", nb.config.WebPort)
			log.Printf("@@@@@ Dial home network: %v addr: %v home: %v\n", network, addr, target)

			return net.Dial(network, target)
		}

		if IsPeer(hostport[0]) {
			log.Printf("@@@@@ Dial peer: %v addr: %v\n", network, addr)

			addr, tld := ConvertTLD(hostport[0])
			//tld := TLD(hostport[0])
			id := ToPeerID(tld)
			if id == "" {
				return nil, fmt.Errorf("Peer invalid: %v", hostport[0])
			}
			target := nb.GetPeerProxy(id)
			if target == "" {
				return nil, fmt.Errorf("Peer not reachable: %v", hostport[0])
			}

			log.Printf("@@@@@ Dial peer: %v addr: %v target: %v\n", network, addr, target)

			dial := proxy.NewConnectDialToProxy(fmt.Sprintf("http://%v", target))
			if dial != nil {
				return dial(network, addr)
			}
			return nil, fmt.Errorf("Peer proxy error: %v", target)
		}

		log.Printf("@@@@@ Dial network: %v addr: %v\n", network, addr)

		return net.Dial(network, addr)
	}

	proxy.Tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            dial,
		DialTLS:         nil,
		Proxy:           p,
	}

	proxy.Verbose = true

	// non proxy request handling
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("@@@@@ NonproxyHandler req: %v\n", req)

		// fwdHandler, _ := forward.New()
		// fwdHost := fmt.Sprintf("127.0.0.1:%v", nb.config.WebPort)
		// req.URL.Host = fwdHost
		// req.Host = fwdHost
		// fwdHandler.ServeHTTP(w, req)

		//TODO check host is in peer id format
		http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", 500)
	})

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("@@@@@ on request Proto: %v method: %v url: %v\n", req.Proto, req.Method, req.URL)
			if req.Method == "CONNECT" {
				panic("boom")
			}
			return req, nil
		})

	// localhost
	var isLocalHost = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			hostport := strings.Split(req.URL.Host, ":")
			b := IsLocalHost(hostport[0])
			log.Printf("@@@@@ isLocalHost: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isLocalHost()).DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if nb.config.Local {
				return req, nil
			}
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusForbidden,
				fmt.Sprintf("Nice try: %v", req.URL.Host))
		})

	// home node - forward to k8s
	var isHome = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			hostport := strings.Split(req.URL.Host, ":")
			b := nb.IsHome(hostport[0])
			log.Printf("@@@@@ isHome: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isHome()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		//req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

		hostport := strings.Split(req.URL.Host, ":")
		addr := nb.ResolveAddr(hostport[0])

		//target := fmt.Sprintf("127.0.0.1:%v", nb.config.WebPort)

		//
		req.Host = addr
		//req.URL.Scheme = "http"
		req.URL.Host = addr
		log.Printf("@@@@@ Home request modified addr: %v req: %v\n", addr, req)

		return req, nil
	})

	// peer node
	var isPeer = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			hostport := strings.Split(req.URL.Host, ":")
			b := nb.IsPeer(hostport[0])
			log.Printf("@@@@@ isPeer: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isPeer()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Peer-ID", nb.My.ID)

		// Resolve peer id
		hostport := strings.Split(req.URL.Host, ":")
		addr := nb.ResolveAddr(hostport[0])

		//
		req.Host = addr
		//req.URL.Scheme = "http"
		req.URL.Host = addr
		log.Printf("@@@@@ Peer request modified addr: %v req: %v\n", addr, req)

		return req, nil
	})

	

	hostport := fmt.Sprintf(":%v", port)
	log.Println("Proxy listening on: " + hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
