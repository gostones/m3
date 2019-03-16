package internal

import (
	"fmt"
	"log"
	"sync"
)

// Peer is
type Peer struct {
	Port    int
	Peer    string
	Latency string
	Muxer   string
	Streams []struct {
		Protocol string
	}

	Rank      int // -1, 0, 1 ...
	timestamp int64
}

// Neighborhood is
type Neighborhood struct {
	Peers  map[string]*Peer
	My     *Node
	Router *RouteRegistry
	// W3ProxyHost string
	config *Config
	min    int
	max    int

	sync.Mutex
}

// NewNeighborhood is
func NewNeighborhood(c *Config) *Neighborhood {
	nb := &Neighborhood{
		Peers:  make(map[string]*Peer, 15),
		config: c,
		min:    0,
		max:    5,
	}

	return nb
}

// GetPeers is
func (r *Neighborhood) GetPeers() []string {
	r.Lock()
	defer r.Unlock()

	addresses := make([]string, 0, len(r.Peers))
	for _, v := range r.Peers {
		if v.Rank > 0 {
			addr := fmt.Sprintf("127.0.0.1:%v", v.Port)
			addresses = append(addresses, addr)
		}
	}
	log.Printf("@@@ addresses: %v\n", addresses)
	return addresses
}

// IsReady tests if node is available
func (r *Neighborhood) IsReady() bool {
	return r.My != nil
}

func (r *Neighborhood) setPeer(p *Peer) {
	r.Lock()
	defer r.Unlock()
	r.Peers[p.Peer] = p
}

func (r *Neighborhood) getPeer(id string) *Peer {
	r.Lock()
	defer r.Unlock()
	p, ok := r.Peers[id]
	if ok {
		return p
	}
	return nil
}

func (r *Neighborhood) AddPeerProxy(id string) string {
	return r.GetPeerTarget(id)
}

// GetPeerTarget returns peer proxy host:port
func (r *Neighborhood) GetPeerTarget(id string) string {
	log.Printf("@@@ GetPeerTarget: id: %v\n", id)

	p := r.getPeer(id)
	if p != nil && p.Port > 0 && p.Rank > 0 {
		addr := fmt.Sprintf("127.0.0.1:%v", p.Port)
		return addr
	}
	//add it
	p = r.addPeer(id)
	addr := fmt.Sprintf("127.0.0.1:%v", p.Port)

	return addr
}

func (r *Neighborhood) addPeer(id string) *Peer {
	if id == r.My.ID {
		//TODO
		panic("attempt to add self as peer: " + id)
	}

	// close old connection
	p := r.getPeer(id)
	if p != nil && p.Port > 0 {
		p2pForwardClose(p.Port, id)
	}

	port := FreePort()
	var err error
	log.Printf("@@@ addPeer: id: %v port: %v\n", id, port)

	err = p2pForward(port, id)

	rank := -1
	if err == nil {
		ok := p2pIsLive(port)
		if ok {
			rank = 1
		}
	}
	log.Printf("@@@ addPeer id: %v port: %v rank: %v err: %v\n", id, port, rank, err)

	p = &Peer{
		Peer:      id,
		Port:      port,
		Rank:      rank,
		timestamp: CurrentTime(),
	}

	// add
	r.setPeer(p)

	return p
}
