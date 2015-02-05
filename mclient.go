/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package main

import (
	"client"
	"fmt"
	"github.com/gaal/go-options/options"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const VERSION = "1.0"

const SPEC = `
Memento - A backup system
Usage: mclient [OPTIONS]
--
h,help                Print this help
v,version             Print version
p,port=               Set port to listen
l,listen=             Set listen address
`

func serve(addr string) {
	var cmd []byte
	var err error

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			// handle error
			fmt.Printf("Error: %v\n", err)
		}

		cmd, err = ioutil.ReadAll(conn)
		client.Parse(cmd, conn)
		conn.Close()
	}
}

func main() {
	var port, listen, address string

	s := options.NewOptions(SPEC)

	// Check if options isn't passed
	if len(os.Args[1:]) <= 0 {
		s.PrintUsageAndExit("No option specified")
	}

	opts := s.Parse(os.Args[1:])

	// Print version and exit
	if opts.GetBool("version") {
		fmt.Println("Memento client " + VERSION)
		os.Exit(0)
	}

	// Print help and exit
	if opts.GetBool("help") {
		s.PrintUsageAndExit("Memento client " + VERSION)
	}

	// Get port to listen
	if opts.GetBool("port") {
		port = opts.Get("port")
	} else {
		fmt.Println("No port specified")
		os.Exit(1)
	}

	// Get address to listen
	if opts.GetBool("listen") {
		if addr := net.ParseIP(opts.Get("listen")); addr != nil {
			listen = addr.String()
		} else {
			log.Fatalln("Invalid IP address")
		}
	}

	if listen == "" {
		address = ":" + port
	} else {
		address = listen + ":" + port
	}

	serve(address)
}
