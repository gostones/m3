package main

import (
	"fmt"
	list "github.com/emirpasic/gods/lists/arraylist"
	"sync"
)

type Backends struct {
	// Length    int
	current   int
	Addresses *list.List

	max int
	sync.Mutex
}

func NewBackends() *Backends {
	addresses := list.New()
	current := 0
	// length := 0

	return &Backends{
		// Length:    length,
		current:   current,
		max:       128,
		Addresses: addresses,
	}
}

func (self *Backends) NextAddress() (addess string) {
	self.Lock()

	index := self.current

	self.current = self.current + 1
	// if self.current > self.Addresses.Size()-1 {
	// 	self.current = 0
	// }
	var addr string
	if self.current >= self.max {
		self.current = 0
		item, _ := self.Addresses.Get(index)
		addr = item.(string)
	} else if self.current > self.Addresses.Size()-1 {
		//TODO ipfs forward
		port := FreePort()
		addr = fmt.Sprintf("localhost:%v", port)
		self.Addresses.Add(addr)
		go forward(addr, "localhost:10080")
	} else {
		item, _ := self.Addresses.Get(index)
		addr = item.(string)
	}

	self.Unlock()
	//

	return addr
}

func (self *Backends) Add(addresses ...string) {
	self.Lock()

	for _, item := range addresses {
		// self.Addresses = append(self.Addresses, item)
		self.Addresses.Add(item)
	}
	// self.Length = len(self.Addresses)

	self.Unlock()
}
