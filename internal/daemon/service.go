package daemon

import (
	"fmt"
	"os"
	// "os/signal"
	// "syscall"

	"github.com/dhnt/m3/internal"
	//"github.com/dhnt/m3/internal/pm"
	"github.com/takama/daemon"
)

const (
	// name of the service
	name        = "dhnt.m3"
	description = "M3 Service"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{}

var stdlog = internal.Stdlog
var errlog = internal.Errlog

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// // Install the service into the system
// func (service *Service) Install(args ...string) (string, error) {
// 	stdlog.Printf("calling super Install os.Args: %v len: %v", os.Args, len(os.Args))

// 	return service.Daemon.Install()
// }

// Remove uninstalls the service and all corresponding files from the system
func (service *Service) Remove() (string, error) {
	stdlog.Printf("calling super Remove os.Args: %v len: %v", os.Args, len(os.Args))
	_, err := service.Daemon.Status()
	if err != nil {
		service.Daemon.Stop()
	}
	return service.Daemon.Remove()
}

// Start the service
func (service *Service) Start() (string, error) {
	stdlog.Printf("calling super Start os.Args: %v len: %v", os.Args, len(os.Args))
	return service.Daemon.Start()
}

// // Stop the service
// func (service *Service) Stop() (string, error) {
// 	stdlog.Printf("calling super Stop os.Args: %v len: %v", os.Args, len(os.Args))
// 	return service.Daemon.Stop()
// }

// // Status - check the service status
// func (service *Service) Status() (string, error) {
// 	stdlog.Printf("calling super status os.Args: %v len: %v", os.Args, len(os.Args))
// 	return service.Daemon.Status()
// }

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	stdlog.Printf("Manage args: %v len: %v", os.Args, len(os.Args))

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

	stdlog.Printf("Manage set up args: %v len: %v", os.Args, len(os.Args))
	// port := internal.GetDaemonPort()
	// pm.StartServer("", port)
	internal.StartGPM()

	return "gpm exited", nil

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// // Set up listener for defined host and port
	// port := internal.GetDaemonPort()
	// listener := pm.NewServer("", port)

	// defer listener.Stop()
	// listener.Start()

	// // loop work cycle with accept connections or interrupt
	// // by system signal
	// for {
	// 	select {
	// 	case killSignal := <-interrupt:
	// 		stdlog.Println("Got signal:", killSignal)
	// 		stdlog.Println("Stoping listening on ", listener.Addr())
	// 		listener.Stop()

	// 		if killSignal == os.Interrupt {
	// 			return "Daemon was interruped by system signal", nil
	// 		}
	// 		return "Daemon was killed", nil
	// 	}
	// }

	// never happen, but need to complete code
	// return usage, nil
}

// func init() {
// 	stdlog = log.New(os.Stdout, "", 0)
// 	errlog = log.New(os.Stderr, "", 0)
// }

// Run daemon service
func Run() {
	stdlog.Printf("Daemon run args: %v", os.Args)
	//internal.DumpEnv()

	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{
		Daemon: srv,
	}

	stdlog.Printf("Calling Manage. service: %v", service)

	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

// func dumpEnv() {
// 	stdlog.Println("dump env ...")

// 	for _, pair := range os.Environ() {
// 		stdlog.Println(pair)
// 	}
// }
