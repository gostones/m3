package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func startServer(listenPort int, backends *Backends) {
	port := strconv.Itoa(listenPort)
	fmt.Println("Starting server on port ", port)

	addr, _ := net.ResolveTCPAddr("tcp", ":"+port)

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Could not listen on port because:", err.Error())
		return
	}

	// listener.SetDeadline(time.Now().Add(time.Second * 10))

	defer func() {
		listener.Close()
		fmt.Println("Listener closed")
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error occurred accepting a connection", err.Error())
			continue
		}

		conn.SetDeadline(time.Now().Add(time.Second * 60))

		addr := backends.NextAddress()
		if addr == "" {
			conn.Close()
			continue
		}
		go handleConnection(conn, addr)
	}
}

func handleConnection(cliConn net.Conn, srvAddr string) {
	srvConn, err := net.Dial("tcp", srvAddr)
	if err != nil {
		fmt.Printf("Could not connect to server (%q), connection dropping\n", srvAddr)
		return
	}

	// close the connections when done
	defer func() {
		srvConn.Close()
		cliConn.Close()
	}()

	go io.Copy(cliConn, srvConn)
	io.Copy(srvConn, cliConn)
}
