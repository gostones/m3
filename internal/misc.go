package internal

import (
	"bufio"
	"fmt"

	"github.com/ilius/crock32"
	"github.com/jpillora/backoff"
	"github.com/multiformats/go-multihash"

	"net"
	"os/exec"
	"path/filepath"
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

// ToPeerID returns b58-encoded ID. it converts to b58 if b32-encoded.
func ToPeerID(s string) string {
	m, err := multihash.FromB58String(s)
	if err == nil {
		return m.B58String()
	}

	c, err := crock32.Decode(s)
	if err != nil {
		return ""
	}

	m, err = multihash.Cast(c)
	if err == nil {
		return m.B58String()
	}

	return ""
}

// ToPeerAddr returns b32-encoded ID. it converts to b32 if B58-encoded.
func ToPeerAddr(s string) string {
	m, err := multihash.FromB58String(s)
	if err == nil {
		return strings.ToLower(crock32.Encode(m))
	}

	//normalize/validate
	d, err := crock32.Decode(s)
	if err == nil {
		m, err = multihash.Cast(d)
		if err != nil {
			return ""
		}
		return strings.ToLower(crock32.Encode(m))
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

// TLD returns last part of a domain name
func TLD(domain string) string {
	sa := strings.Split(domain, ".")
	s := sa[len(sa)-1]

	return s
}

// PeerTLD splits and returns peer id after strippig off m3
func PeerTLD(domain string) string {
	sa := strings.Split(domain, ".")
	s := sa[len(sa)-1]
	if s == "m3" {
		s = sa[len(sa)-2]
	}
	return s
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

func ToTimestamp(d time.Time) int64 {
	return d.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Execute sets up env and runs file
func Execute(base, file string, args ...string) error {
	// binary, err := exec.LookPath(file)
	// if err != nil {
	// 	return err
	// }
	binary := filepath.Join(base, file)
	cmd := exec.Command(binary, args...)

	//
	cmd.Env = DefaultEnviron(base)

	//
	// cmdOut, err := cmd.StdoutPipe()
	// cmdErr, _ := cmd.StderrPipe()

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		logger.Println("error creating stdout", err)
		return err
	}

	scanner := bufio.NewScanner(cmdOut)
	go func() {
		for scanner.Scan() {
			logger.Println(">", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		logger.Println("error starting cmd", err)
		return err
	}

	return cmd.Wait()
}
