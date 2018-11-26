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

	b = IsPeerID("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
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

	b = IsPeerAddress("http://QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ/")
	if !b {
		t.Fail()
	}
}

func TestPeerIDHex(t *testing.T) {
	s := PeerIDHex("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
	fmt.Println(s)
	if s == "" {
		t.Fail()
	}
}
