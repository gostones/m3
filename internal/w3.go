package internal

import (
	"encoding/json"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"time"
)

type Health struct {
	Healthy   bool  `json:"healthy"`
	Timestamp int64 `json:"timestamp"`
}

func toTimestamp(d time.Time) int64 {
	return d.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

//W3Proxy start a local proxy to w3
func W3Proxy(port int) {
	hostport := fmt.Sprintf(":%v", port)
	proxy := goproxy.NewProxyHttpServer()
	proxy.ConnectDial = nil

	//proxy.Tr.Proxy = nil
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		m := &Health{
			Healthy:   true,
			Timestamp: toTimestamp(time.Now()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(m)
		fmt.Fprintf(w, string(b))
	})
	proxy.Verbose = true
	log.Printf("local proxy listening on: %v\n", hostport)
	log.Fatal(http.ListenAndServe(hostport, proxy))
}
