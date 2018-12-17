// https://github.com/elazarl/goproxy
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func httpproxy(port int, nb *Neighborhood) {

	log.Printf("@@@ ProxyURL: %v\n", nb.config.ProxyURL)

	//
	proxy := goproxy.NewProxyHttpServer()

	p := func(req *http.Request) (*url.URL, error) {
		log.Printf("@@@ Proxying: %v %v %v url: %v\n", req.Proto, req.Method, req.Host, req.URL)

		hostport := strings.Split(req.URL.Host, ":")
		proxyURL := nb.config.ProxyURL

		if IsLocalHost(hostport[0]) || IsHome(hostport[0]) {
			return nil, nil
		}

		if IsPeer(hostport[0]) {
			log.Printf("@@@ Proxy url: %v\n", req.URL)

			tld := TLD(hostport[0])
			id := ToPeerID(tld)
			if id == "" {
				return nil, fmt.Errorf("Peer invalid: %v", hostport[0])
			}
			target := nb.GetPeerProxy(id)
			if target == "" {
				return nil, fmt.Errorf("Peer not reachable: %v", hostport[0])
			}

			log.Printf("@@@ Proxy peer url: %v target: %v\n", req.URL, target)

			proxyURL, _ = url.Parse(fmt.Sprintf("http://%v", target))
			return proxyURL, nil
		}

		return proxyURL, nil
	}

	dial := func(network, addr string) (net.Conn, error) {
		hostport := strings.Split(addr, ":")

		if IsHome(hostport[0]) {
			target := fmt.Sprintf("127.0.0.1:%v", nb.config.WebPort)
			log.Printf("@@@ Dial home network: %v addr: %v home: %v\n", network, addr, target)

			return net.Dial(network, target)
		}

		if IsPeer(hostport[0]) {
			log.Printf("@@@ Dial peer: %v addr: %v\n", network, addr)

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

			log.Printf("@@@ Dial peer: %v addr: %v target: %v\n", network, addr, target)

			dial := proxy.NewConnectDialToProxy(fmt.Sprintf("http://%v", target))
			if dial != nil {
				return dial(network, addr)
			}
			return nil, fmt.Errorf("Peer proxy error: %v", target)
		}

		log.Printf("@@@ Dial network: %v addr: %v\n", network, addr)

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
		log.Printf("@@@ NonproxyHandler req: %v\n", req)

		//TODO check host is in peer id format
		http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", 500)
	})

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("@@@ OnRequest Proto: %v method: %v url: %v\n", req.Proto, req.Method, req.URL)
			log.Printf("@@@ OnRequest request: %v\n", req)

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
			log.Printf("@@@ isLocalHost: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	var isBlocked = func(host string) bool {
		hostport := strings.Split(host, ":")
		port := "80"
		if len(hostport) > 1 {
			port = hostport[1]
		}
		for _, v := range nb.config.Blocked {
			if v == port {
				return true
			}
		}
		return false
	}
	proxy.OnRequest(isLocalHost()).DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if nb.config.Local {
				if isBlocked(req.URL.Host) {
					//silently ignore by returning empty OK response
					log.Printf("@@@ OnRequest blocking host: %v url: %v\n", req.Host, req.URL)

					return req, goproxy.NewResponse(req,
						goproxy.ContentTypeText, http.StatusOK, "")
				}
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
			log.Printf("@@@ isHome: %v host: %v\n", b, req.URL.Host)
			return b
		}
	}
	proxy.OnRequest(isHome()).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Peer-ID", nb.My.ID)

		hostport := strings.Split(req.URL.Host, ":")
		addr := nb.ResolveAddr(hostport[0])
		//
		req.Host = addr
		//req.URL.Scheme = "http"
		req.URL.Host = addr
		log.Printf("@@@ OnRequest home modified url: %v header: %v\n", req.URL, req.Header)

		return req, nil
	})

	// peer node
	var isPeer = func() goproxy.ReqConditionFunc {
		return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			hostport := strings.Split(req.URL.Host, ":")
			b := nb.IsPeer(hostport[0])
			log.Printf("@@@ isPeer: %v host: %v\n", b, req.URL.Host)
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
		log.Printf("@@@ OnRequest peer modified url: %v header: %v\n", req.URL, req.Header)

		return req, nil
	})

	// response
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if r != nil {
			log.Printf("@@@ OnResponse status: %v length: %v\n", r.StatusCode, r.ContentLength)
		}

		log.Printf("@@@ OnResponse response: %v\n", r)

		return r
	})

	log.Printf("Proxy listening on: %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), proxy))
}
