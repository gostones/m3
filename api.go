package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"log"
	"net/url"
)

var client = resty.New()

const (
	protocol_name = "/x/www/1.0"
)

// ipfs p2p listen /x/kickass/1.0 /ip4/127.0.0.1/tcp/$APP_PORT
func p2pListen(appPort int) error {
	target := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", appPort)

	resp, err := client.R().
		SetMultiValueQueryParams(url.Values{
			"arg": []string{protocol_name, target},
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
			"arg": []string{protocol_name, listen, target},
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/p2p/forward")

	log.Printf("Status: %v\n", resp.Status())
	log.Println(resp)

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
