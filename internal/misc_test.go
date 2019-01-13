package internal

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

	b = IsPeer("home")
	if b {
		t.Fail()
	}

	b = IsPeer("example.com")
	if b {
		t.Fail()
	}

	//peer
	b = IsPeer("92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8")
	if !b {
		t.Fail()
	}

	b = IsPeer("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
	if !b {
		t.Fail()
	}

	b = IsPeer("QmTFdcQY12fjxv6kELzQA4zXBxiva8xcunrmTYZto8DFUk")
	if !b {
		t.Fail()
	}
}

func TestToPeerAddr(t *testing.T) {
	s := ToPeerAddr("localhost")
	if s != "" {
		t.Fail()
	}

	s = ToPeerAddr("home")
	if s != "" {
		t.Fail()
	}

	s = ToPeerAddr("example.com")
	if s != "" {
		t.Fail()
	}

	s = ToPeerAddr("92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8")
	if s != "92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8" {
		t.Fail()
	}

	s = ToPeerAddr("QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ")
	if s != "92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8" {
		t.Fail()
	}

	s = ToPeerAddr("QmTFdcQY12fjxv6kELzQA4zXBxiva8xcunrmTYZto8DFUk")
	if s != "920j8197p3mq6fdb1xhpserz22x146pzmnr2hj987ay8nea39fr40h" {
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

	s = ToPeerID("92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8")
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

func TestIsLocalHost(t *testing.T) {
	s := "example.com"
	b := IsLocalHost(s)
	if b {
		t.Fail()
	}

	s = "home"
	b = IsLocalHost(s)
	if b {
		t.Fail()
	}

	s = "example.home"
	b = IsLocalHost(s)
	if b {
		t.Fail()
	}

	//
	s = "localhost"
	b = IsLocalHost(s)
	if !b {
		t.Fail()
	}

	s = "127.0.0.1"
	b = IsLocalHost(s)
	if !b {
		t.Fail()
	}
}

func TestIsHOme(t *testing.T) {
	s := "example.com"
	b := IsHome(s)
	if b {
		t.Fail()
	}

	s = "localhost"
	b = IsHome(s)
	if b {
		t.Fail()
	}

	s = "127.0.0.1"
	b = IsHome(s)
	if b {
		t.Fail()
	}

	//
	s = "home"
	b = IsHome(s)
	if !b {
		t.Fail()
	}

	s = "any.home"
	b = IsHome(s)
	if !b {
		t.Fail()
	}
}

// func TestSetDefaultEnv(t *testing.T) {
// 	SetDefaultEnv()
// 	for _, nv := range os.Environ() {
// 		t.Log(nv)
// 	}
// }
