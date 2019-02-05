package internal

// partially based on https://github.com/google/tcpproxy/blob/master/cmd/tlsrouter/config.go

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Backend struct {
	Hostname string
	Port     int
	Healthy  bool
}

// A Route maps a match on a domain name to a backend.
type Route struct {
	re      *regexp.Regexp
	pattern string
	Backend []*Backend
	Proxy   bool
}

// RouteRegistry stores the routing configuration.
type RouteRegistry struct {
	mu     sync.Mutex
	MyID   string
	MyAddr string
	Routes []*Route
}

// func (r *RouteRegistry) SetDefault(target string) {
// 	r.defaultRoute = &Route{
// 		Domain: "",
// 		Backend: []*Backend{
// 			&Backend{
// 				Host:    target,
// 				Healthy: true,
// 			},
// 		},
// 	}
// }

// func (r *RouteRegistry) GetDefault() *Route {
// 	return r.DefaultRoute
// }

// func (r *RouteRegistry) Set(name string, target *Route) {
// 	r.route[name] = target
// }

// func (r *RouteRegistry) Delete(name string) {
// 	delete(r.route, name)
// }

// func (r *RouteRegistry) Add(name, target string) {

// 	route, found := r.route[name]
// 	if found {
// 		route.Backend = append(route.Backend, &Backend{
// 			Host:    target,
// 			Healthy: true,
// 		})
// 		return
// 	}

// 	r.route[name] = &Route{
// 		Domain: name,
// 		Backend: []*Backend{
// 			&Backend{
// 				Host:    target,
// 				Healthy: true,
// 			},
// 		},
// 	}
// }

// //TODO cache result
// func (r *RouteRegistry) Match(name string) *Route {
// 	var matched []*Route
// 	domain := "." + name
// 	for _, v := range r.route {
// 		if name == v.Domain {
// 			return v
// 		}
// 		if strings.HasSuffix(domain, v.Domain) && strings.HasPrefix(v.Domain, ".") {
// 			matched = append(matched, v)
// 		}
// 	}

// 	//
// 	if len(matched) == 0 {
// 		return r.defaultRoute
// 	}
// 	//find longest match
// 	found := matched[0]
// 	for i := 0; i < len(matched); i++ {
// 		if len(matched[i].Domain) > len(found.Domain) {
// 			found = matched[i]
// 		}
// 	}
// 	return found
// }

func (c *RouteRegistry) parseBackend(s string) *Backend {
	var be Backend
	hp := strings.Split(s, ":")
	be.Hostname = hp[0]
	if len(hp) > 1 {
		if p, err := strconv.Atoi(hp[1]); err == nil {
			be.Port = p
		}
	}
	return &be
}

// expand variables: myid
func (c *RouteRegistry) expandVar(s string) string {
	mapper := func(n string) string {
		switch n {
		case "myid":
			return c.MyAddr
		}
		return ""
	}
	return os.Expand(s, mapper)
}

func (c *RouteRegistry) parseDomain(s string) (*regexp.Regexp, string, error) {
	s = c.expandVar(s)

	if len(s) >= 2 && s[0] == '/' && s[len(s)-1] == '/' {
		re, err := regexp.Compile(s[1 : len(s)-1])
		return re, "", err
	}
	return nil, s, nil
}

// Match returns the backend for hostname, and whether to use proxy.
func (c *RouteRegistry) Match(hostname string) ([]*Backend, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, r := range c.Routes {
		if r.re != nil && r.re.MatchString(hostname) {
			return r.Backend, r.Proxy
		}
		if r.pattern != "" {
			if matched, err := filepath.Match(r.pattern, hostname); matched && err == nil {
				return r.Backend, r.Proxy
			}
		}
	}
	return nil, false
}

// Read replaces current config
func (c *RouteRegistry) Read(reader io.Reader) error {
	var routes []*Route

	s := bufio.NewScanner(reader)
	for s.Scan() {
		if strings.HasPrefix(strings.TrimSpace(s.Text()), "#") {
			// Comment, ignore.
			continue
		}

		fs := strings.Fields(s.Text())
		switch len(fs) {
		case 0:
			continue
		case 1:
			return fmt.Errorf("invalid entry: %q", s.Text())
		case 2:
			re, pa, err := c.parseDomain(fs[0])
			if err != nil {
				return err
			}

			routes = append(routes, &Route{
				re:      re,
				pattern: pa,
				Backend: []*Backend{c.parseBackend(fs[1])},
				Proxy:   false,
			})
		case 3:
			re, pa, err := c.parseDomain(fs[0])
			if err != nil {
				return err
			}
			if strings.ToLower(fs[2]) != "proxy" {
				return errors.New("invalid proxy flag")
			}
			routes = append(routes, &Route{
				re:      re,
				pattern: pa,
				Backend: []*Backend{c.parseBackend(fs[1])},
				Proxy:   true,
			})
		default:
			// TODO: multiple backends?
			return fmt.Errorf("multiple backeds not supported yet: %v", s.Text())
		}
	}
	if err := s.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.Routes = routes

	return nil
}

// ReadFile replaces the current routes with one read from path.
func (c *RouteRegistry) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return c.Read(f)
}

// ReadString replaces the current routes with one read from cfg.
func (c *RouteRegistry) ReadString(cfg string) error {
	b := bytes.NewBufferString(cfg)
	return c.Read(b)
}

// NewRouteRegistry instantiates a new route registry
func NewRouteRegistry(myid string) *RouteRegistry {
	addr := ToPeerAddr(myid)
	return &RouteRegistry{
		MyID:   myid,
		MyAddr: addr,
	}
}
