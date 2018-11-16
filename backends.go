package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Backends struct {
	Length    int
	current   int
	Addresses []string

	sync.Mutex
}

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

func (self *Backends) NextAddress() (addess string) {
	self.Lock()

	index := self.current

	self.current = self.current + 1
	if self.current > self.Length-1 {
		self.current = 0
	}

	self.Unlock()
	return self.Addresses[index]
}

func (self *Backends) Add(addresses ...string) {
	self.Lock()

	for _, item := range addresses {
		self.Addresses = append(self.Addresses, item)
	}
	self.Length = len(self.Addresses)

	self.Unlock()
}

type Peer struct {
	Address   string
	Rank      int
	timestamp int64
}

type Neighborhood struct {
	Peers map[string]*Peer
	min   int
	max   int

	sync.Mutex
}

func NewNeighborhood() *Neighborhood {
	min := 2

	n := &Neighborhood{
		Peers: make(map[string]*Peer, 15),
	}
	n.InitPeers(min)
	return n
}

func (r *Neighborhood) InitPeers(cnt int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < cnt; i++ {
		id := fmt.Sprintf("peer_%v", rnd.Intn(36))
		r.AddPeer(id)
	}
}

func (r *Neighborhood) GetPeers() []string {
	addresses := make([]string, len(r.Peers))

	i := 0
	for _, v := range r.Peers {
		addresses[i] = v.Address
		i++
	}
	return addresses
}

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
