package pm

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/dhnt/m3/internal"
)

var logger = internal.Stdlog

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
	gpm     *internal.GPM
	mux     sync.Mutex
}

func (r *Handler) start() {
	r.mux.Lock()
	defer r.mux.Unlock()
	if !r.running {
		r.running = true
		r.gpm.Start()
	}
}

func (r *Handler) stop() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.gpm.Stop()

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

// Stop gracefully terminates the server listener.
func (r *Server) Stop() (err error) {
	if r.listener != nil {
		err = r.listener.Close()
	}
	return
}

// Addr returns host:port
func (r *Server) Addr() string {
	addr := fmt.Sprintf("%v:%v", r.Host, r.Port)
	return addr
}

// Start Runs the RPC server.
func (r *Server) Start() (err error) {
	go r.Run()
	return nil
}

// Run initializes the RPC server.
func (r *Server) Run() (err error) {
	logger.Println("RPC Serve starting ...")

	internal.DumpEnv()

	handler := &Handler{
		gpm: internal.NewGPM(),
	}
	err = rpc.Register(handler)
	if err != nil {
		return
	}
	r.listener, err = net.Listen("tcp", r.Addr())
	if err != nil {
		return
	}

	defer handler.stop()
	handler.start()

	logger.Println("RPC Serve accepting ...")

	rpc.Accept(r.listener)

	logger.Println("RPC Serve exited.")

	return
}

// NewServer creates a RPC server
func NewServer(port int) *Server {
	return &Server{
		Host: "",
		Port: port,
	}
}

func StartServer(host string, port int) {
	s := Server{
		Host: host,
		Port: port,
	}

	defer s.Stop()

	logger.Printf("starting: %v\n", s)

	s.Run()

	logger.Printf("exited: %v\n", s)
}
