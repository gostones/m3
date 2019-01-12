package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	m3 "github.com/dhnt/m3/internal"
	"github.com/parnurzeal/gorequest"
)

type message struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	port := flag.Int("port", m3.GetIntEnv("PORT", 18082), "port")
	svc := flag.Int("service", m3.GetIntEnv("M3_PORT", 18080), "M3 service port")

	flag.Parse()

	gpm := m3.NewGPM()
	//
	endpoint := fmt.Sprintf("http://localhost:%v", *svc)
	address := fmt.Sprintf(":%v", *port)

	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		log.Println("status")

		w.Header().Set("Content-Type", "application/json")

		err := ping(endpoint)
		status := "running"
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			status = fmt.Sprintf("%v", err)
		}
		response(status, w)
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Starting M3 ...")
		go gpm.Start()

		response("started", w)
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Stopping M3 ...")
		gpm.Stop()

		response("stopped", w)
	})

	log.Printf("Listening: %v\n", address)

	log.Fatal(http.ListenAndServe(address, nil))
}

func ping(u string) error {
	request := gorequest.New()
	resp, _, errs := request.
		Head(u).
		End()

	log.Printf("response: %v err %v\n", resp, errs)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func response(s string, w http.ResponseWriter) {
	now := m3.ToTimestamp(time.Now())
	m := message{
		Status:    s,
		Timestamp: now,
	}
	json.NewEncoder(w).Encode(m)
}
