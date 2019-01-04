package internal

import (
	"testing"
)

func TestRouteRegistry(t *testing.T) {
	reg := NewRouteRegistry()

	reg.SetDefault("localhost:80")

	reg.Add("home", "localhost:8080")
	reg.Add(".home", "localhost:8081")
	reg.Add("matrix.home", "localhost:8082")
	reg.Add(".matrix.home", "localhost:8083")
	reg.Add("atrix.home", "localhost:8084")

	//
	r := reg.GetDefault()
	if r.Backend[0].Host != "localhost:80" {
		t.Fail()
	}

	r = reg.Match("home")
	if r.Backend[0].Host != "localhost:8080" {
		t.Fail()
	}
	r = reg.Match("riot.home")
	if r.Backend[0].Host != "localhost:8081" {
		t.Fail()
	}
	r = reg.Match("matrix.home")
	if r.Backend[0].Host != "localhost:8082" {
		t.Fail()
	}
	r = reg.Match("federation.matrix.home")
	if r.Backend[0].Host != "localhost:8083" {
		t.Fail()
	}
	r = reg.Match("atrix.home")
	if r.Backend[0].Host != "localhost:8084" {
		t.Fail()
	}
}
