package internal

import (
	"fmt"
	"log"
	"net/http"
	//"net/url"
	//"os"
	"testing"
)

func webserver(port int) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", r.URL)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func TestIsLive(t *testing.T) {
	//t.Skip()

	port := FreePort()

	var cfg = &Config{}
	//cfg.WebHost = fmt.Sprintf("http://localhost:%v", FreePort())

	proxyPort := FreePort()
	cfg.Pals = []string{""}
	//cfg.ProxyURL, _ = url.Parse(os.Getenv("http_proxy"))

	t.Logf("Configuration: %v\n", cfg)

	nb := NewNeighborhood(cfg)

	addr := fmt.Sprintf("127.0.0.1:%v", port)
	target := fmt.Sprintf("127.0.0.1:%v", proxyPort)
	go HTTPProxy(proxyPort, nb)
	go Forward(addr, target)
	go webserver(FreePort())
	t.Logf("addr: %v target: %v", addr, target)

	ok := p2pIsLive(port)

	if !ok {
		t.Fail()
	}

	ok = p2pIsProxy(port)

	if !ok {
		t.Fail()
	}
}

func TestIsP2pProxy(t *testing.T) {
	//t.Skip()

	id := "QmTFdcQY12fjxv6kELzQA4zXBxiva8xcunrmTYZto8DFUk"
	//id := "QmXG428k4Aa6Fchp7buub2pK4Fa2nbhcTfznL7oVSGWRRZ"
	//
	port := FreePort()
	var cfg = &Config{}
	//cfg.WebPort = 5001
	//proxyPort = FreePort()
	cfg.Pals = []string{id}

	t.Logf("Configuration: %v\n", cfg)

	err := p2pForward(port, id)
	if err != nil {
		t.Fail()
	}

	ok := p2pIsProxy(port)

	if !ok {
		t.Fail()
	}
}

func TestP2pCloseAll(t *testing.T) {
	t.Skip()

	err := P2PCloseAll()
	if err != nil {
		t.Fail()
	}
}
