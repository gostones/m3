package internal

import (
	"fmt"
	"net/url"
)

// Config is application settings
type Config struct {
	Port     int
	Home     string
	WebProxy *url.URL
	Local    bool
	Blocked  []string
	Proxy    []string
	Alias    map[string]string
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
