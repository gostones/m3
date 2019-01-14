package pm

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// Client contains the configuration options for a RPC client.
type Client struct {
	Host string
	Port int

	client *rpc.Client
}

func NewClient(host string, port int) (*Client, error) {
	c := Client{
		Host: host,
		Port: port,
	}

	err := c.Connect()

	return &c, err
}

// Connect initializes the underlying RPC client.
func (r *Client) Connect() error {
	addr := fmt.Sprintf("%v:%v", r.Host, r.Port)

	// r.client, err = rpc.Dial("tcp", addr)
	timeout := 5 * time.Second
	client, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}

	r.client = rpc.NewClient(client)
	return nil
}

// Close gracefully terminates the underlying client.
func (r *Client) Close() error {
	if r.client != nil {
		err := r.client.Close()
		return err
	}

	return nil
}

// Start requests service start.
func (r *Client) Start() (*Response, error) {
	var (
		req = &Request{}
		res = new(Response)
	)

	err := r.client.Call(handlerStart, req, res)
	return res, err
}

// Stop requests service stop.
func (r *Client) Stop() (*Response, error) {
	var (
		req = &Request{}
		res = new(Response)
	)

	err := r.client.Call(handlerStop, req, res)
	return res, err
}

// Status gets service status.
func (r *Client) Status() (*Response, error) {
	var (
		req = &Request{}
		res = new(Response)
	)

	err := r.client.Call(handlerStatus, req, res)
	return res, err
}
