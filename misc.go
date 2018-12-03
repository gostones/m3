package main

import (
	"fmt"
	"github.com/jpillora/backoff"
	"github.com/multiformats/go-multihash"
	"net"
	"net/url"
	"strconv"
	"time"
)

// BackoffDuration is
func BackoffDuration() func(error) {
	b := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    15 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	return func(rc error) {
		secs := b.Duration()

		fmt.Printf("rc: %v sleeping %v\n", rc, secs)
		time.Sleep(secs)
		if secs.Nanoseconds() >= b.Max.Nanoseconds() {
			b.Reset()
		}
	}
}

// FreePort is
func FreePort() int {
	l, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// CurrentTime is
func CurrentTime() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// IsPeerID checks if the string s is a valid hex encoded peer address
func IsPeerID(s string) bool {
	_, err := multihash.FromHexString(s)
	return err == nil
}

// IsPeerAddress checks if string s is a valid Url with host being a valid peer ID
func IsPeerAddress(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	scheme := u.Scheme
	host := u.Hostname()
	return (scheme == "http" || scheme == "https") && IsPeerID(host)
}

// PeerIDHex converts B58-encoded multihash peer ID to hex-encoded string
func PeerIDHex(s string) string {
	h, err := multihash.FromB58String(s)
	if err == nil {
		return h.HexString()
	}
	return ""
}

// PeerIDB58 converts hex-encoded multihash peer ID to B58-encoded string
func PeerIDB58(s string) string {
	h, err := multihash.FromHexString(s)
	if err == nil {
		return h.B58String()
	}
	return ""
}

// ParseInt parses s into int
func ParseInt(s string, v int) int {
	if s == "" {
		return v
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		i = v
	}
	return i
}
