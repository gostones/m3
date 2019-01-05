package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gostones/lib"
	"github.com/parnurzeal/gorequest"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Health struct {
	Healthy   bool  `json:"healthy"`
	Timestamp int64 `json:"timestamp"`
}

func HealthHandlerFunc(proxyURL string) http.HandlerFunc {
	const elapse int64 = 60000 //one min
	last := toTimestamp(time.Now())
	healthy := false

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
