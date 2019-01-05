package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Health struct {
	Healthy   bool  `json:"healthy"`
	Timestamp int64 `json:"timestamp"`
}

func HealthHandler(w http.ResponseWriter, req *http.Request) {
	m := &Health{
		Healthy:   true,
		Timestamp: toTimestamp(time.Now()),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(m)
	fmt.Fprintf(w, string(b))
}

func toTimestamp(d time.Time) int64 {
	return d.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
