package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gostones/lib"
	"github.com/parnurzeal/gorequest"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Health struct {
	Healthy   bool  `json:"healthy"`
	Timestamp int64 `json:"timestamp"`
}

func HealthHandlerFunc(proxyURL string) http.HandlerFunc {
	const elapse int64 = 60000 //one min
	last := ToTimestamp(time.Now())
	healthy := false
	mutex := &sync.Mutex{}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		now := toTimestamp(time.Now())
		if !healthy || now-last > elapse {
			healthy = pingProxy(proxyURL)
			last = now
		}
		m := &Health{
			Healthy:   healthy,
			Timestamp: now,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(m)
		fmt.Fprintf(w, string(b))
	})
}

func toTimestamp(d time.Time) int64 {
	return d.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func pingProxy(proxy string) bool {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	testSite := []string{
		"https://www.google.com/",
		"https://aws.amazon.com/",
		"https://azure.microsoft.com/",
	}

	request := gorequest.New().Proxy(proxy)

	//
	err := lib.Retry(func() error {
		idx := rnd.Intn(len(testSite))
		resp, _, errs := request.
			Head(testSite[idx]).
			End()

		log.Printf("Proxy: %v response: %v err %v\n", proxy, resp, errs)
		if len(errs) > 0 {
			return errs[0]
		}
		return nil
	})

	return err == nil
}

const pacFile = `
function FindProxyForURL(url, host) {
	return "PROXY %v";
}
`

// PACHandlerFunc handles PAC file request
func PACHandlerFunc(proxyURL string) http.HandlerFunc {
	URL, _ := url.Parse(proxyURL)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		s := fmt.Sprintf(pacFile, URL.Host)
		w.Write([]byte(s))
	})
}

// MuxHandlerFunc multiplexes requests
func MuxHandlerFunc(proxyURL string) http.HandlerFunc {
	mux := http.NewServeMux()
	mux.HandleFunc("/proxy.pac", PACHandlerFunc(proxyURL))
	mux.HandleFunc("/health", HealthHandlerFunc(proxyURL))
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req)
	})
}
