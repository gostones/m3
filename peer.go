package main

import (
	"fmt"
	"sync"
)

// // Peer is
// type Peer struct {
// 	Address   string
// 	Rank      int // -1, 0, 1 ...
// 	timestamp int64
// }

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
			addresses = append(addresses, v.Addr)
		}
	}
	fmt.Printf("@@@@ addresses: %v\n", addresses)
	return addresses
}

// Monitor is
func (r *Neighborhood) Monitor() {
	job := func() {
		// clean up stale connections
		for k, v := range r.Peers {
			if v.Rank == 0 {
				delete(r.Peers, k)
			}
		}

		//
		for i := 0; i < r.current; i++ {
			// TODO
			id := config.Node.ID
			peer, found := r.Peers[id]

			fmt.Printf("@@@@ Peer ID: %v found: %v\n", id, found)

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
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	// target := fmt.Sprintf("localhost:%v", config.ProxyPort)
	// go forward(addr, target)

	p2pForward(port, id)

	fmt.Printf("@@@@ Added peer: %v addr: %v\n", id, addr)

	r.Peers[id] = &Peer{
		Addr:      addr,
		Rank:      1,
		timestamp: CurrentTime(),
	}
}
