package internal

import (
	"fmt"
	"log"
	"strings"
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
	Peers       map[string]*Peer
	My          *Node
	Home        *RouteRegistry
	W3ProxyHost string
	config      *Config
	min         int
	max         int

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

	//
	// const concurrency = 8 // max
	// ch := make(chan string, concurrency)

	// go func() {
	// 	for c := range ch {
	// 		go func(id string) {
	// 			p := nb.checkPeer(id)
	// 			if p.Rank > 0 {
	// 				nb.setPeer(p)
	// 			}
	// 		}(c)
	// 	}
	// }()

	nb.init()

	//go nb.addSelf()
	//go nb.addPals(ch)
	//go nb.Monitor()
	return nb
}

func (r *Neighborhood) init() {
	node, err := p2pID()
	if err != nil {
		panic(err)
	}
	r.My = &node
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

// // Monitor is
// func (r *Neighborhood) Monitor() {
// 	job := func() {
// 		if !r.IsReady() {
// 			return
// 		}

// 		// clean up stale connections
// 		// for k, v := range r.Peers {
// 		// 	if v.Rank == 0 {
// 		// 		delete(r.Peers, k)
// 		// 	}
// 		// }

// 		cur := 0
// 		for _, v := range r.Peers {
// 			if v.Rank > 0 {
// 				cur++
// 			}
// 		}
// 		if cur >= r.max {
// 			log.Printf("@@@ current count: %v max: %v, no new peers will be added\n", cur, r.max)
// 			return
// 		}

// 		//
// 		peers, err := p2pPeers()
// 		if err != nil {
// 			log.Printf("@@@ get peers: %v\n", err)
// 			return
// 		}

// 		cnt := len(peers)
// 		log.Printf("@@@ get peers, count: %v\n", cnt)
// 		if cnt <= 0 {
// 			return
// 		}

// 		const concurrency = 32 // max
// 		ch := make(chan string, concurrency)
// 		go func() {
// 			for i := 0; i < cnt; i++ {
// 				p := peers[i]
// 				id := p.Peer
// 				peer, found := r.Peers[id] // TODO?

// 				log.Printf("@@@ Peer ID: %v found: %v\n", id, found)

// 				if found {
// 					peer.Rank++
// 				} else {
// 					ch <- id
// 				}
// 			}
// 		}()

// 		for c := range ch {
// 			go func(id string) {
// 				p := r.checkPeer(id)
// 				r.setPeer(p)
// 			}(c)
// 		}
// 	}

// 	Every(15).Seconds().Run(job)
// }

// func (r *Neighborhood) addSelf() {
// 	node, err := p2pID()
// 	if err != nil {
// 		panic(err)
// 	}
// 	r.My = &node

// 	//
// 	p := r.checkPeer(r.My.ID)
// 	r.setPeer(p)
// }

// IsReady tests if node is available
func (r *Neighborhood) IsReady() bool {
	return r.My != nil
}

func (r *Neighborhood) resolveAlias(s string) string {
	a, err := Alias(s)
	if err == nil && a != "" {
		alias, b := r.config.Alias[a]
		if b {
			return alias
		}
	}
	return s
}

// IsHome checks if host is not localhost but a home node
func (r *Neighborhood) IsHome(host string) bool {
	host = r.ResolveAddr(host)
	return IsHome(host)
}

// IsPeer checks if host is not local host or home and is a valid remote peer
func (r *Neighborhood) IsPeer(host string) bool {
	host = r.ResolveAddr(host)
	return IsPeer(host)
}

// ResolveAddr resolves s into canonical address form
// localhost
// *.home
// *.<hex> -- peer node
// *.<tld> -- world wide web address
func (r *Neighborhood) ResolveAddr(s string) string {
	// TODO alias name service
	s = r.resolveAlias(s)

	//localhost
	if IsLocalHost(s) {
		return s
	}

	//*.home
	if IsHome(s) {
		return s
	}

	//*.<hex>
	sa := strings.Split(s, ".")
	le := len(sa) - 1
	id := ToPeerID(sa[le])
	if id != "" {
		if id == r.My.ID {
			id = "home"
		}
		if le > 0 {
			return strings.Join(sa[:le], ".") + "." + id
		}
		return id
	}

	//*.<tld> internet address
	return s
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

// func (r *Neighborhood) addPals(ch chan<- string) {

// 	job := func() {
// 		if !r.IsReady() {
// 			return
// 		}

// 		pals := r.config.Pals
// 		cnt := len(pals)
// 		log.Printf("@@@ pals count: %v\n", cnt)
// 		if cnt <= 0 {
// 			return
// 		}

// 		for i := 0; i < cnt; i++ {
// 			id := pals[i]
// 			peer, found := r.Peers[id] // TODO?

// 			log.Printf("@@@ Peer ID: %v found: %v\n", id, found)

// 			if found {
// 				peer.Rank++
// 			} else {
// 				ch <- id
// 			}
// 		}
// 	}

// 	Every(1).Minutes().Run(job)
// }

func (r *Neighborhood) AddPeerProxy(id string) string {
	return r.GetPeerTarget(id)
}

// GetPeerTarget returns peer proxy host:port
func (r *Neighborhood) GetPeerTarget(id string) string {
	log.Printf("@@@ GetPeerTarget: id: %v\n", id)

	// if id == r.My.ID {
	// 	//TODO
	// 	return ""
	// 	//return fmt.Sprintf("localhost:%v", r.config.WebPort)
	// }

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

// func (r *Neighborhood) checkPeer(id string) *Peer {
// 	port := FreePort()
// 	addr := fmt.Sprintf("127.0.0.1:%v", port)

// 	var err error
// 	self := (id == r.My.ID)

// 	log.Printf("@@@ checkPeer: id: %v self: %v addr: %v\n", id, self, addr)

// 	if self {
// 		target := fmt.Sprintf("127.0.0.1:%v", r.config.ProxyPort)
// 		go Forward(addr, target)
// 	} else {
// 		err = p2pForward(port, id)
// 	}
// 	rank := -1
// 	if err == nil {
// 		ok := p2pIsLive(port)
// 		if ok {
// 			rank = 1
// 		} else if !self {
// 			p2pForwardClose(port, id) // no www support
// 		}
// 	}
// 	log.Printf("@@@ check peer: self: %v ID: %v addr: %v rank: %v err: %v\n", self, id, addr, rank, err)

// 	return &Peer{
// 		Peer:      id,
// 		Port:      port,
// 		Rank:      rank,
// 		timestamp: CurrentTime(),
// 	}
// }

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
