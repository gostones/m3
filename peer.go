package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Peer is
type Peer struct {
	Address   string
	Rank      int // -1, 0, 1 ...
	timestamp int64
}

// Neighborhood is
type Neighborhood struct {
	Peers   map[string]*Peer
	current int
	min     int
	max     int

	sync.Mutex
}

// NewNeighborhood is
func NewNeighborhood() *Neighborhood {
	nb := &Neighborhood{
		Peers:   make(map[string]*Peer, 15),
		current: 5,
		min:     0,
		max:     8,
	}
	nb.Monitor()
	return nb
}

// GetPeers is
func (r *Neighborhood) GetPeers() []string {
	r.Lock()
	defer r.Unlock()

	addresses := make([]string, 0, len(r.Peers))
	for _, v := range r.Peers {
		if v.Rank > 0 {
			addresses = append(addresses, v.Address)
		}
	}
	fmt.Printf("@@@@ addresses: %v\n", addresses)
	return addresses
}

// Monitor is
func (r *Neighborhood) Monitor() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	job := func() {
		// clean up stale connections
		for k, v := range r.Peers {
			if v.Rank == 0 {
				delete(r.Peers, k)
			}
		}

		//
		for i := 0; i < r.current; i++ {
			n := rnd.Intn(36)
			id := fmt.Sprintf("peer_%v", n)
			peer, found := r.Peers[id]

			fmt.Printf("@@@@ Rand: %v found: %v\n", n, found)

			if found {
				peer.Rank++
			} else {
				r.addPeer(id)
			}
		}
	}

	Every(1).Minutes().Run(job)
}

// addPeer is
func (r *Neighborhood) addPeer(id string) {
	port := FreePort()
	addr := fmt.Sprintf("localhost:%v", port)

	go forward(addr, "localhost:10080")

	fmt.Printf("@@@@ Added peer: %v addr: %v\n", id, addr)

	r.Peers[id] = &Peer{
		Address:   addr,
		Rank:      1,
		timestamp: CurrentTime(),
	}
}
