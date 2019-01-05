package internal

import (
	"fmt"
	//"net/url"
	"strings"
)

// Config is application settings
type Config struct {
	Port    int
	Local   bool
	Blocked []string
	Home    []string
	Web     []string
	Alias   map[string]string
}

// ListFlags is for collecting an array of command line arguments
type ListFlags []string

func (r *ListFlags) String() string {
	return fmt.Sprintf("%v", *r)
}

// Set appends the value
func (r *ListFlags) Set(value string) error {
	*r = append(*r, value)
	return nil
}

type Backend struct {
	Host    string // name:port
	Healthy bool
}

type Route struct {
	Domain  string
	Backend []*Backend
}

type RouteRegistry struct {
	route        map[string]*Route // domain -> route
	defaultRoute *Route
}

func (r *RouteRegistry) SetDefault(target string) {
	r.defaultRoute = &Route{
		Domain: "",
		Backend: []*Backend{
			&Backend{
				Host:    target,
				Healthy: true,
			},
		},
	}
}

func (r *RouteRegistry) GetDefault() *Route {
	return r.defaultRoute
}

func (r *RouteRegistry) Set(name string, target *Route) {
	r.route[name] = target
}

func (r *RouteRegistry) Delete(name string) {
	delete(r.route, name)
}

func (r *RouteRegistry) Add(name, target string) {

	route, found := r.route[name]
	if found {
		route.Backend = append(route.Backend, &Backend{
			Host:    target,
			Healthy: true,
		})
		return
	}

	r.route[name] = &Route{
		Domain: name,
		Backend: []*Backend{
			&Backend{
				Host:    target,
				Healthy: true,
			},
		},
	}
}

//TODO cache result
func (r *RouteRegistry) Match(name string) *Route {
	var matched []*Route
	domain := "." + name
	for _, v := range r.route {
		if name == v.Domain {
			return v
		}
		if strings.HasSuffix(domain, v.Domain) && strings.HasPrefix(v.Domain, ".") {
			matched = append(matched, v)
		}
	}

	//
	if len(matched) == 0 {
		return r.defaultRoute
	}
	//find longest match
	found := matched[0]
	for i := 0; i < len(matched); i++ {
		if len(matched[i].Domain) > len(found.Domain) {
			found = matched[i]
		}
	}
	return found
}

func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		route: make(map[string]*Route),
	}
}
