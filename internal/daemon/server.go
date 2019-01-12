package daemon

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/dhnt/m3/internal"
)

type Request struct {
}

type Response struct {
	Message string
	Status  bool
}

const (
	handlerStart  = "Handler.Start"
	handlerStop   = "Handler.Stop"
	handlerStatus = "Handler.Status"
)

// Handler holds the methods to be exposed by the RPC
// server as well as properties
type Handler struct {
	running bool
	mux     sync.Mutex
}

func (r *Handler) start() {
	r.mux.Lock()
	defer r.mux.Unlock()
	if !r.running {
		r.running = true
		internal.StartGPM()
	}
}

func (r *Handler) stop() {
	r.mux.Lock()
	defer r.mux.Unlock()
	internal.StopGPM()

	r.running = false
}

// Start starts service
func (r *Handler) Start(req Request, res *Response) error {
	r.start()
	res.Status = true
	res.Message = "started"

	return nil
}

// Stop stops service
func (r *Handler) Stop(req Request, res *Response) error {
	r.stop()
	res.Status = true
	res.Message = "stopped"

	return nil
}

// Status returns service status
func (r *Handler) Status(req Request, res *Response) error {
	res.Status = r.running
	res.Message = fmt.Sprintf("running - %v", r.running)

	return nil
}

// Server holds the configuration used to initiate
// an RPC server.
type Server struct {
	Host string
	Port int

	listener net.Listener
}

// Close gracefully terminates the server listener.
func (r *Server) Close() (err error) {
	if r.listener != nil {
		err = r.listener.Close()
	}
	return
}

// Addr
func (r *Server) Addr() string {
	addr := fmt.Sprintf("%v:%v", r.Host, r.Port)
	return addr
}

// Serve initializes the RPC server.
func (r *Server) Serve() (err error) {
	handler := &Handler{}
	err = rpc.Register(handler)
	if err != nil {
		return
	}
	r.listener, err = net.Listen("tcp", r.Addr())
	if err != nil {
		return
	}

	handler.start()
	defer handler.stop()

	rpc.Accept(r.listener)

	return
}

// func StartServer(host string, port int) {
// 	s := Server{
// 		Host: host,
// 		Port: port,
// 	}

// 	defer s.Close()

// 	log.Printf("starting: %v\n", s)

// 	s.Serve()

// 	log.Printf("exited: %v\n", s)
// }
