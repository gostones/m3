package pm

import (
	"fmt"
	"github.com/dhnt/m3/internal"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var logger = internal.Logger()

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
	base string
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

	//internal.DumpEnv()

	handler := &Handler{
		gpm: internal.NewGPM(r.base),
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

	done := make(chan error, 1)
	go func() {
		rpc.Accept(r.listener)
		done <- fmt.Errorf("rpc accept exited")
	}()

	signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case err = <-done:
		logger.Println("error:", err)
	case <-signalChan:
		// case <-r.signalChan:
	}

	logger.Println("RPC Serve exited.")

	return
}

// NewServer creates rpc server
func NewServer(base string, host string, port int) *Server {
	return &Server{
		base: base,
		Host: host,
		Port: port,
	}
}

// StartServer runs pm server
func StartServer(base, host string, port int) {

	s := NewServer(base, host, port)

	defer s.Stop()

	logger.Printf("starting: %v %v %v", base, host, port)

	s.Run()

	logger.Println("exited")
}
