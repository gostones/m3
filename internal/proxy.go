// https://github.com/elazarl/goproxy
package internal

import (
	"crypto/tls"
	"fmt"
	"github.com/elazarl/goproxy"
	//"github.com/elazarl/goproxy/ext/auth"
	"github.com/dhnt/m3/internal/lb"

	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func redirectHost(r *http.Request, host, body string) *http.Response {
	resp := &http.Response{}
	resp.Request = r
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", "text/plain")

	u := *r.URL
	u.Host = host
	resp.Header.Set("Location", u.String())

	resp.StatusCode = http.StatusMovedPermanently
	resp.Status = http.StatusText(resp.StatusCode)
	buf := bytes.NewBufferString(body)
	resp.ContentLength = int64(buf.Len())
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

// HTTPProxy dispatches request based on network addr
func HTTPProxy(port int, nb *Neighborhood) {

	//log.Printf("@@@ ProxyURL: %v\n", nb.config.ProxyURL)

	//
	proxy := goproxy.NewProxyHttpServer()

	p := func(req *http.Request) (*url.URL, error) {
		log.Printf("@@@ Proxy: %v %v %v url: %v\n", req.Proto, req.Method, req.Host, req.URL)

		hostport := strings.Split(req.URL.Host, ":")
		hostport[0] = nb.ResolveAddr(hostport[0])

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
			target := nb.GetPeerTarget(id)
			if target == "" {
				return nil, fmt.Errorf("Peer not reachable: %v", hostport[0])
			}

			log.Printf("@@@ Proxy peer url: %v target: %v\n", req.URL, target)

			proxyURL, _ := url.Parse(fmt.Sprintf("http://%v", target))
			return proxyURL, nil
		}

		return nb.config.WebProxy, nil
	}

	dial := func(network, addr string) (net.Conn, error) {
		hostport := strings.Split(addr, ":")
		hostport[0] = nb.ResolveAddr(hostport[0])

		if IsHome(hostport[0]) {
			target := nb.config.Home
			log.Printf("@@@ Dial home network: %v addr: %v home: %v\n", network, addr, target)

			return net.Dial(network, target)
		}

		if IsPeer(hostport[0]) {
			log.Printf("@@@ Dial peer: %v addr: %v\n", network, addr)

			addr, tld := ConvertTLD(hostport[0])
			id := ToPeerID(tld)
			if id == "" {
				return nil, fmt.Errorf("Peer invalid: %v", hostport[0])
			}
			target := nb.GetPeerTarget(id)
			if target == "" {
				return nil, fmt.Errorf("Peer not reachable: %v", hostport[0])
			}

			log.Printf("@@@ Dial peer: %v addr: %v target: %v\n", network, addr, target)

			dial := proxy.NewConnectDialToProxyWithHandler(fmt.Sprintf("http://%v", target), func(req *http.Request) {
				log.Printf("\n@@@ Dial NewConnectDialToProxyWithHandler peer: %v addr: %v target: %v\nreq: %v\n", network, addr, target, req)
			})
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
	proxy.ConnectDial = nil
	//
	proxy.Verbose = true

	// non proxy request handling
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("@@@ NonproxyHandler req: %v\n", req)

		//TODO check host is in peer id format
		http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", 500)
	})

	// auth
	// auth.ProxyBasic(proxy, "m3_realm", func(user, passwd string) bool {
	// 	//TODO verify peer is who it claims to be
	// 	//user is the peer id and pwd is: peer_addr,timestamp
	// 	//after decrypting with peer's public key
	// 	//return user == json[0]
	// 	return true
	// })

	// request
	var isBlocked = func(port string) bool {
		for _, v := range nb.config.Blocked {
			if v == port {
				return true
			}
		}
		return false
	}

	var convertAliasTLD = func(host string) (string, bool) {
		sa := strings.Split(host, ".")
		le := len(sa)
		tld := sa[le-1]
		//
		alias, ok := nb.config.Alias[tld]
		if !ok {
			return "", false
		}
		addr := ToPeerAddr(alias)
		if addr == "" {
			return "", false
		}
		sa[le-1] = addr
		return strings.Join(sa, "."), true
	}

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			log.Printf("\n\n\n##################\n")

			log.Printf("@@@ OnRequest Proto: %v method: %v url: %v\n", req.Proto, req.Method, req.URL)
			log.Printf("@@@ OnRequest request: %v\n", req)

			//
			hostport := strings.Split(req.URL.Host, ":")

			// block localhost or specified local ports
			if IsLocalHost(hostport[0]) {
				log.Printf("@@@ OnRequest isLocal: %v\n", req.URL.Host)
				if nb.config.Local {
					port := "80"
					if len(hostport) > 1 {
						port = hostport[1]
					}
					if isBlocked(port) {
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
			}

			// alias
			if h, err := Alias(hostport[0]); err == nil {
				addr, ok := convertAliasTLD(h)
				if !ok {
					return req, goproxy.NewResponse(req,
						goproxy.ContentTypeText, http.StatusNotFound,
						fmt.Sprintf("Alias invalid: %v", hostport[0]))
				}
				return req, redirectHost(req, addr, fmt.Sprintf("Redirect: %v", hostport[0]))
			}

			return req, nil
		})

	// response
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		log.Printf("\n--------------------\n")

		if r != nil {
			r.Header.Add("X-Peer-Id", nb.My.ID)
			log.Printf("@@@ OnResponse status: %v length: %v\n", r.StatusCode, r.ContentLength)
		}

		log.Printf("@@@ OnResponse response: %v\n", r)

		return r
	})

	log.Printf("Proxy listening on: %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), proxy))
}

// StartProxy starts proxy services
func StartProxy(cfg *Config) {
	// clean up old p2p connections
	err := P2PCloseAll()
	if err != nil {
		panic(err)
	}
	//
	log.Printf("Configuration: %v\n", cfg)

	nb := NewNeighborhood(cfg)

	// local web proxy
	lbPort := FreePort()
	cfg.WebProxy, _ = url.Parse(fmt.Sprintf("http://127.0.0.1:%v", lbPort))

	localProxyPort := FreePort()
	go LocalProxy(localProxyPort)

	//TODO dynamic proxy
	backends := []string{fmt.Sprintf("localhost:%v", localProxyPort)}

	for _, v := range cfg.Proxy {
		addr := nb.AddPeerProxy(v)
		backends = append(backends, addr)
	}

	//
	port := cfg.Port

	log.Printf("web proxy load balancer: %v backends: %v\n", lbPort, backends)
	log.Printf("proxy/p2p port: %v\n", port)

	go lb.Start(lbPort, backends, true)

	P2PListen(port)
	HTTPProxy(port, nb)
}
