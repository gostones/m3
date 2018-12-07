package main

import (
	"fmt"
	"testing"
)

func TestIsLive(t *testing.T) {
	t.Skip()

	var cfg = &Config{}
	cfg.Port = FreePort()
	cfg.WebPort = 5001
	cfg.ProxyPort = FreePort()
	cfg.Pals = []string{"QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ"}

	t.Logf("Configuration: %v\n", cfg)

	nb := NewNeighborhood(cfg)

	addr := fmt.Sprintf("127.0.0.1:%v", cfg.Port)
	target := fmt.Sprintf("127.0.0.1:%v", cfg.ProxyPort)
	go httpproxy(cfg.ProxyPort, nb)
	go forward(addr, target)

	t.Logf("addr: %v target: %v", addr, target)

	ok := p2pIsLive(cfg.Port)

	if !ok {
		t.Fail()
	}

	ok = p2pIsProxy(cfg.Port)

	if !ok {
		t.Fail()
	}
}

func TestIsP2pProxy(t *testing.T) {
	//t.Skip()

	id := "QmTFdcQY12fjxv6kELzQA4zXBxiva8xcunrmTYZto8DFUk"
	//id := "QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ"
	//
	var cfg = &Config{}
	cfg.Port = FreePort()
	cfg.WebPort = 5001
	cfg.ProxyPort = FreePort()
	cfg.Pals = []string{id}

	t.Logf("Configuration: %v\n", cfg)

	err := p2pForward(cfg.Port, id)
	if err != nil {
		t.Fail()
	}

	ok := p2pIsProxy(cfg.Port)

	if !ok {
		t.Fail()
	}
}

func TestP2pCloseAll(t *testing.T) {
	t.Skip()

	err := p2pCloseAll()
	if err != nil {
		t.Fail()
	}
}
