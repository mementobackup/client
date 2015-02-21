/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package main

import (
	"client"
	"code.google.com/p/goconf/conf"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"github.com/gaal/go-options/options"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
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
s,ssl=                Set SSL config file
`

func tlsserve(addr, key, private string) net.Listener {
	var err error
	var ln net.Listener

	cert, err := tls.LoadX509KeyPair(key, private)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}

	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader

	ln, err = tls.Listen("tcp", addr, &config)

	return ln
}

func plainserve(addr string) net.Listener {
	var err error
	var ln net.Listener

	ln, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	return ln
}

func serve(addr string, ssl *conf.ConfigFile) {
	var cmd []byte
	var ln net.Listener

	if ssl != nil {
		key, _ := ssl.GetString("ssl", "key")
		private, _ := ssl.GetString("ssl", "private")
		ln = tlsserve(addr, key, private)
	} else {
		ln = plainserve(addr)
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

	if opts.GetBool("ssl") {
		cfg, err := conf.ReadConfigFile(opts.Get("ssl"))
		if err != nil {
			// handle error
			log.Fatalf("Error: %v\n", err)
		}

		serve(address, cfg)
	} else {
		serve(address, nil)
	}
}
