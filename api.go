package main

import (
	"encoding/json"
	"fmt"
	"github.com/gostones/lib"
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

// Peers is
type Peers struct {
	Peers []Peer
}

// Node is
type Node struct {
	ID              string
	PublicKey       string
	Addresses       []string
	AgentVersion    string
	ProtocolVersion string
}

// ipfs p2p listen /x/www/1.0 /ip4/127.0.0.1/tcp/$APP_PORT
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

// ipfs p2p forward /x/www/1.0 /ip4/127.0.0.1/tcp/$SOME_PORT /ipfs/$SERVER_ID
func p2pForward(port int, serverID string) error {
	listen := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", port)
	target := fmt.Sprintf("/ipfs/%v", serverID)

	log.Printf("p2pForward %v %v\n", listen, target)

	resp, err := client.R().
		SetMultiValueQueryParams(url.Values{
			"arg": []string{protocolWWW, listen, target},
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/p2p/forward")

	log.Printf("p2pForward  %v %v response: %v err: %v\n", listen, target, resp, err)

	return err
}

func p2pForwardClose(port int, serverID string) error {
	listen := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", port)
	target := fmt.Sprintf("/ipfs/%v", serverID)

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"protocol":       protocolWWW,
			"listen-address": listen,
			"target-address": target,
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/p2p/close")

	log.Printf("close forward  %v %v response: %v err: %v\n", listen, target, resp, err)

	return err
}

func p2pCloseAll() error {
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"protocol": protocolWWW,
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:5001/api/v0/p2p/close")

	log.Printf("close all response: %v err: %v\n", resp, err)

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

func p2pIsLive(port int) bool {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{
		"http://home/",
	}
	proxy := fmt.Sprintf("http://127.0.0.1:%v", port)
	request := gorequest.New().Proxy(proxy)

	//
	err := lib.Retry(func() error {
		idx := rnd.Intn(len(tests))
		resp, _, errs := request.
			Head(tests[idx]).
			End()

		log.Printf("proxy: %v response: %v err %v\n", proxy, resp, errs)
		if len(errs) > 0 {
			return errs[0]
		}
		return nil
	})

	return err == nil
}

func p2pIsProxy(port int) bool {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{
		"https://www.google.com/",
		"https://aws.amazon.com/",
		"https://azure.microsoft.com/",
	}
	proxy := fmt.Sprintf("http://127.0.0.1:%v", port)
	request := gorequest.New().Proxy(proxy)

	//
	err := lib.Retry(func() error {
		idx := rnd.Intn(len(tests))
		resp, _, errs := request.
			Head(tests[idx]).
			End()

		log.Printf("Proxy: %v response: %v err %v\n", proxy, resp, errs)
		if len(errs) > 0 {
			return errs[0]
		}
		return nil
	})

	return err == nil
}
