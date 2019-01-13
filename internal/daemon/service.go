package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhnt/m3/internal"
	"github.com/takama/daemon"
)

const (
	// name of the service
	name        = "dhnt.m3"
	description = "M3 Service"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	stdlog.Printf("Manage args: %v", os.Args)

	usage := "Usage: m3d install | uninstall | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "uninstall":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	port := internal.GetDaemonPort()

	listener := &Server{
		Port: port,
	}

	stdlog.Printf("Starting gpm daemon")

	defer listener.Close()

	go listener.Serve()

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())

			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

// Run daemon service
func Run() {
	stdlog.Println("Daemon ...")
	dumpEnv()

	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}

	stdlog.Printf("Calling Manage. service: %v", service)

	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

func dumpEnv() {
	stdlog.Println("dump env ...")

	for _, pair := range os.Environ() {
		stdlog.Println(pair)
	}
}
