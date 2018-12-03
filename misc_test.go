package main

import (
	"fmt"
	"testing"
)

func TestBackoffDuration(t *testing.T) {
	bo := BackoffDuration()
	count := 0
	rc := func() error {
		count++
		return fmt.Errorf("count: %v", count)
	}()
	bo(rc)
}

func TestIsPeerID(t *testing.T) {
	b := IsPeerID("localhost")
	if b {
		t.Fail()
	}

	b = IsPeerID("www.google.com")
	if b {
		t.Fail()
	}

	b = IsPeerID("1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08")
	if !b {
		t.Fail()
	}
}

func TestIsPeerAddress(t *testing.T) {
	b := IsPeerAddress("http://localhost:5001/")
	if b {
		t.Fail()
	}

	b = IsPeerAddress("https://www.google.com")
	if b {
		t.Fail()
	}

	b = IsPeerAddress("http://1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08/")
	if !b {
		t.Fail()
	}
}

func TestPeerIDHex(t *testing.T) {
	s := PeerIDHex("QmTFdcQY12fjxv6kELzQA4zXBxiva8xcunrmTYZto8DFUk")
	if s != "1220490149ec3a5ccf6ac3d8db2ec7c42e8486b7e95c0a324a0eaf22ae50d2fc1011" {
		t.Fail()
	}
}

func TestPeerIDB58(t *testing.T) {
	s := PeerIDB58("1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08")
	if s != "QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ" {
		t.Fail()
	}
}
