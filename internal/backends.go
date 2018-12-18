package internal

import (
	"sync"
)

// Backends is
type Backends struct {
	Length    int
	current   int
	Addresses []string

	nb *Neighborhood
	sync.Mutex
}

// NewBackends is
func NewBackends(nb *Neighborhood) *Backends {
	current := 0
	length := 0

	return &Backends{
		Length:  length,
		current: current,
		nb:      nb,
	}
}

// NextAddress is
func (r *Backends) NextAddress() (addess string) {
	r.Lock()

	index := r.current

	r.current = r.current + 1
	if r.current > r.Length-1 {
		r.current = 0
		r.Addresses = r.nb.GetPeers()
	}

	r.Unlock()
	if index < len(r.Addresses) {
		return r.Addresses[index]
	}
	return ""
}

// Add is
func (r *Backends) Add(addresses ...string) {
	r.Lock()

	for _, item := range addresses {
		r.Addresses = append(r.Addresses, item)
	}
	r.Length = len(r.Addresses)

	r.Unlock()
}
