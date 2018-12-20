package main

import (
	"flag"
	"fmt"
	misc "github.com/dhnt/m3/internal"
	"os"
)

const man = `
NAME
	hexid -- Peer ID conversion

SYNOPSIS
	hexid [--id peerid ] [--hex addr]

DESCRIPTION
	The hexid utility converts peer ID to address or vice versa.

	--id peerid
            Print peer hex encoded address of the given peerid
		  
	--hex addr
            Print peer ID of the given peer address a.k.a hex encoded peer ID

`

func usage() {
	fmt.Println(man)
	os.Exit(1)
}

func main() {
	id := flag.String("id", "", "Peer ID in b58 encoding")
	addr := flag.String("hex", "", "Hex encoded Peer ID")

	flag.Parse()

	if *id == "" && *addr == "" {
		usage()
	}

	if *id != "" {
		s := misc.ToPeerAddr(*id)
		if s == "" {
			fmt.Println("Invalid peer ID")
			os.Exit(1)
		}
		fmt.Println(s)
	}

	if *addr != "" {
		s := misc.ToPeerID(*addr)
		if s == "" {
			fmt.Println("Invalid peer address")
			os.Exit(1)
		}
		fmt.Println(s)
	}
	os.Exit(0)
}
