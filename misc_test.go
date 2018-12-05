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

func TestIsPeer(t *testing.T) {
	b := IsPeer("localhost")
	if b {
		t.Fail()
	}

	b = IsPeer("www.google.com")
	if b {
		t.Fail()
	}

	b = IsPeer("1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08")
	if !b {
		t.Fail()
	}

	b = IsPeer("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
	if !b {
		t.Fail()
	}
}

func TestToPeerAddr(t *testing.T) {
	s := ToPeerAddr("localhost")
	if s != "" {
		t.Fail()
	}

	s = ToPeerAddr("example.com")
	if s != "" {
		t.Fail()
	}

	s = ToPeerAddr("1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08")
	if s != "1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08" {
		t.Fail()
	}

	s = ToPeerAddr("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
	if s != "1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08" {
		t.Fail()
	}
}

func TestToPeerID(t *testing.T) {
	s := ToPeerID("localhost")
	if s != "" {
		t.Fail()
	}

	s = ToPeerID("example.com")
	if s != "" {
		t.Fail()
	}

	s = ToPeerID("1220848ba2cbc954d17fc1758a4dc06ec128b21c6ecc1dcfcbdc284809f4a922ba08")
	if s != "QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ" {
		t.Fail()
	}

	s = ToPeerID("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
	if s != "QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ" {
		t.Fail()
	}
}

func TestTLD(t *testing.T) {
	s := "example.com"
	tld := TLD(s)
	if tld != "com" {
		t.Fail()
	}
}

func TestAlias(t *testing.T) {
	_, err := Alias("example.com")
	if err == nil {
		t.Fail()
	}
	a, _ := Alias("my.friend.a")
	if a != "my.friend" {
		t.Fail()
	}
	a, err = Alias("a")
	if a != "" || err != nil {
		t.Fail()
	}
	a, err = Alias(".a")
	if a != "" || err != nil {
		t.Fail()
	}
}
