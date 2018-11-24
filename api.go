package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/resty.v1"
	"log"
	"math/rand"
	"net/url"
	"time"
)

var client = resty.New()

const (
	protocolWWW = "/x/www/1.0"
)

// ipfs p2p listen /x/kickass/1.0 /ip4/127.0.0.1/tcp/$APP_PORT
func p2pListen(appPort int) error {
	target := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", appPort)

	resp, err := client.R().
		SetMultiValueQueryParams(url.Values{
			"arg": []string{protocolWWW, target},
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/p2p/listen")

	log.Printf("Status: %v\n", resp.Status())
	log.Println(resp)

	return err
}

// ipfs p2p forward /x/kickass/1.0 /ip4/127.0.0.1/tcp/$SOME_PORT /ipfs/$SERVER_ID
func p2pForward(port int, serverID string) error {
	listen := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", port)
	target := fmt.Sprintf("/ipfs/%v", serverID)

	resp, err := client.R().
		SetMultiValueQueryParams(url.Values{
			"arg": []string{protocolWWW, listen, target},
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/p2p/forward")

	log.Printf("forward  %v %v response: %v err: %v\n", listen, target, resp, err)

	return err
}

//
// {
//     "Peers": [
//         {
//             "Addr": "<string>"
//             "Peer": "<string>"
//             "Latency": "<string>"
//             "Muxer": "<string>"
//             "Streams": [
//                 {
//                     "Protocol": "<string>"
//                 }
//             ]
//         }
//     ]
// }

// Peer is
type Peer struct {
	Addr    string
	Peer    string
	Latency string
	Muxer   string
	Streams []struct {
		Protocol string
	}

	Rank      int // -1, 0, 1 ...
	timestamp int64
}

// Peers is
type Peers struct {
	Peers []Peer
}

func p2pPeers() ([]Peer, error) {
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/swarm/peers?verbose=true&streams=true&latency=true")

	log.Printf("Status: %v\n", resp.Status())

	p := Peers{}

	json.Unmarshal([]byte(resp.Body()), &p)

	return p.Peers, err
}

//
// {
//     "ID": "<string>"
//     "PublicKey": "<string>"
//     "Addresses": [
//         "<string>"
//     ]
//     "AgentVersion": "<string>"
//     "ProtocolVersion": "<string>"
// }

// Node is
type Node struct {
	ID              string
	PublicKey       string
	Addresses       []string
	AgentVersion    string
	ProtocolVersion string
}

func p2pID() (Node, error) {
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/id")

	log.Printf("Status: %v\n", resp.Status())

	n := Node{}

	json.Unmarshal([]byte(resp.Body()), &n)

	return n, err
}

func p2pIsValid(port int) bool {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{
		"https://www.google.com",
	}
	idx := rnd.Intn(len(tests))

	proxy := fmt.Sprintf("http://127.0.0.1:%v", port)

	request := gorequest.New().Proxy(proxy)
	resp, _, err := request.
		Head(tests[idx]).
		Retry(3, 5*time.Second).
		End()

	log.Printf("Proxy: %v response: %v err %v\n", proxy, resp, err)

	return err == nil
}
