package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Backends is
type Backends struct {
	Length    int
	current   int
	Addresses []string

	sync.Mutex
}

// NewBackends is
func NewBackends() *Backends {
	nb := NewNeighborhood()
	addresses := nb.GetPeers()

	current := 0
	length := 0

	return &Backends{
		Length:    length,
		current:   current,
		Addresses: addresses,
	}
}

// NextAddress is
func (r *Backends) NextAddress() (addess string) {
	r.Lock()

	index := r.current

	r.current = r.current + 1
	if r.current > r.Length-1 {
		r.current = 0
	}

	r.Unlock()
	return r.Addresses[index]
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

// Peer is
type Peer struct {
	Address   string
	Rank      int
	timestamp int64
}

// Neighborhood is
type Neighborhood struct {
	Peers map[string]*Peer
	min   int
	max   int

	sync.Mutex
}

// NewNeighborhood is
func NewNeighborhood() *Neighborhood {
	min := 2

	n := &Neighborhood{
		Peers: make(map[string]*Peer, 15),
	}
	n.InitPeers(min)
	return n
}

// InitPeers is
func (r *Neighborhood) InitPeers(cnt int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < cnt; i++ {
		id := fmt.Sprintf("peer_%v", rnd.Intn(36))
		r.AddPeer(id)
	}
}

// GetPeers is
func (r *Neighborhood) GetPeers() []string {
	addresses := make([]string, len(r.Peers))

	i := 0
	for _, v := range r.Peers {
		addresses[i] = v.Address
		i++
	}
	return addresses
}

// Monitor is
func (r *Neighborhood) Monitor() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	job := func() {
		i := rnd.Intn(36)
		id := fmt.Sprintf("peer_%v", i)
		peer, found := r.Peers[id]
		fmt.Printf("@@@@ Rand: %v found: %v\n", i, found)

		if found {
			peer.Rank = 0
		} else {
			r.AddPeer(id)
		}

	}

	Every(1).Minutes().Run(job)
}

// AddPeer is
func (r *Neighborhood) AddPeer(id string) {

	port := FreePort()
	addr := fmt.Sprintf("localhost:%v", port)
	go forward(addr, "localhost:10080")

	fmt.Printf("@@@@ Added peer: %v addr: %v\n", id, addr)
	r.Peers[id] = &Peer{
		Address:   addr,
		Rank:      1,
		timestamp: currentTime(),
	}
}

func currentTime() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
