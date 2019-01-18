package pm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dhnt/m3/internal"
	"github.com/parnurzeal/gorequest"
)

type message struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type HTTPServer struct {
	gpm *internal.GPM

	endpoint string
	address  string
}

func NewHTTPServer(base, ep string) *HTTPServer {
	return &HTTPServer{
		gpm:      internal.NewGPM(base),
		endpoint: ep,
	}
}

func (r *HTTPServer) Status() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Println("status")

		w.Header().Set("Content-Type", "application/json")

		err := ping(r.endpoint)
		status := "running"
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			status = fmt.Sprintf("%v", err)
		}
		response(status, w)
	}
}

func (r *HTTPServer) Start() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Println("Starting M3 ...")
		r.gpm.Start()

		response("started", w)
	}
}

func (r *HTTPServer) Stop() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Println("Stopping M3 ...")
		r.gpm.Stop()

		response("stopped", w)
	}
}

func StartHTTPServer(base string, port, m3port int) {
	endpoint := fmt.Sprintf("http://localhost:%v", m3port)
	address := fmt.Sprintf(":%v", port)
	s := NewHTTPServer(base, endpoint)

	http.HandleFunc("/status", s.Status())
	http.HandleFunc("/start", s.Start())
	http.HandleFunc("/stop", s.Stop())

	logger.Printf("Listening: %v\n", address)

	logger.Fatal(http.ListenAndServe(address, nil))
}

func ping(u string) error {
	request := gorequest.New()
	resp, _, errs := request.
		Head(u).
		End()

	logger.Printf("response: %v err %v\n", resp, errs)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func response(s string, w http.ResponseWriter) {
	now := internal.ToTimestamp(time.Now())
	m := message{
		Status:    s,
		Timestamp: now,
	}
	json.NewEncoder(w).Encode(m)
}
