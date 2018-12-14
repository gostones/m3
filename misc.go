package main

import (
	"fmt"
	"github.com/jpillora/backoff"
	"github.com/multiformats/go-multihash"
	"net"
	"regexp"
	"strconv"
	"strings"
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

// IsPeer checks if the string s ends in a valid hex-encoded peer address or b58-encoded ID
func IsPeer(s string) bool {
	sa := strings.Split(s, ".")
	le := len(sa) - 1
	id := ToPeerID(sa[le])
	return id != ""
}

// // IsPeerAddress checks if string s is a valid Url with host being a valid peer ID
// func IsPeerAddress(s string) bool {
// 	u, err := url.Parse(s)
// 	if err != nil {
// 		return false
// 	}
// 	scheme := u.Scheme
// 	host := u.Hostname()
// 	return (scheme == "http" || scheme == "https") && IsPeerID(host)
// }

// ToPeerID returns b58-encoded ID. it converts to b58 if hex-encoded.
func ToPeerID(s string) string {
	m, err := multihash.FromB58String(s)
	if err == nil {
		return m.B58String()
	}
	m, err = multihash.FromHexString(s)
	if err == nil {
		return m.B58String()
	}
	return ""
}

// ToPeerAddr returns hex-encoded ID. it converts to hex if B58-encoded.
func ToPeerAddr(s string) string {
	m, err := multihash.FromB58String(s)
	if err == nil {
		return m.HexString()
	}
	m, err = multihash.FromHexString(s)
	if err == nil {
		return m.HexString()
	}
	return ""
}

// // ToPeerIDHex converts B58-encoded multihash peer ID to hex-encoded string
// func ToPeerIDHex(s string) string {
// 	h, err := multihash.FromB58String(s)
// 	if err == nil {
// 		return h.HexString()
// 	}
// 	return ""
// }

// // ToPeerIDB58 converts hex-encoded multihash peer ID to B58-encoded string
// func ToPeerIDB58(s string) string {
// 	h, err := multihash.FromHexString(s)
// 	if err == nil {
// 		return h.B58String()
// 	}
// 	return ""
// }

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

// TLD returns last part of a domain name
func TLD(name string) string {
	sa := strings.Split(name, ".")
	s := sa[len(sa)-1]

	return s
}

// Alias returns the second last part of a domain name ending in .a
// or error if not an alias
func Alias(name string) (string, error) {
	if name != "a" && !strings.HasSuffix(name, ".a") {
		return "", fmt.Errorf("Not an alias: %v", name)
	}
	sa := strings.Split(name, ".")
	if len(sa) == 1 {
		return "", nil
	}
	s := sa[0 : len(sa)-1]

	return strings.Join(s, "."), nil
}

var localHostIpv4RE = regexp.MustCompile(`127\.0\.0\.\d+`)

// IsLocalHost checks whether host is explicitly local host
// taken from goproxy
func IsLocalHost(host string) bool {
	return host == "::1" ||
		host == "0:0:0:0:0:0:0:1" ||
		localHostIpv4RE.MatchString(host) ||
		host == "localhost"
}

var localHomeRE = regexp.MustCompile(`.*\.?home`)

// IsHome checks whether host is explicitly local home node
func IsHome(host string) bool {
	return host == "home" ||
		localHomeRE.MatchString(host)
}
