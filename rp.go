package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	hostProxy map[string]*httputil.ReverseProxy
)

type baseHandle struct{}

func (h baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if fn, ok := hostProxy[host]; ok {
		fn.ServeHTTP(w, r)
		return
	}

	target := "http://" + host

	log.Println("target: ", target)

	remote, err := url.Parse(target)
	if err != nil {
		log.Println("target parse fail:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	hostProxy[host] = proxy
	proxy.ServeHTTP(w, r)
	return

}

func reverseProxy(port int) {
	hostProxy = make(map[string]*httputil.ReverseProxy)
	h := &baseHandle{}
	http.Handle("/", h)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: h,
	}

	//server.ListenAndServeTLS( "server.pem", "server.key")

	log.Fatal(server.ListenAndServe())
}
