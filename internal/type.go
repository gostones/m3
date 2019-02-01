package internal

import (
	"fmt"
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
