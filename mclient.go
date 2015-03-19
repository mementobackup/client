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
	"github.com/go-ini/ini"
	"github.com/op/go-logging"
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
s,ssl=                Set SSL config file
D,debug               Enable debug messages
`

func setLog(debug bool) *logging.Logger {
    var log = logging.MustGetLogger("Memento Client")

    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backendLeveled := logging.AddModuleLevel(backend)

    if debug {
        backendLeveled.SetLevel(logging.DEBUG, "")
    } else {
        backendLeveled.SetLevel(logging.CRITICAL, "")
    }

    logging.SetBackend(backendLeveled)

    return log
}

func main() {
    var log *logging.Logger
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

    // Enable debug
    if opts.GetBool("debug") {
        log = setLog(true)
    } else {
        log = setLog(false)
    }

	// Get port to listen
	if opts.GetBool("port") {
		port = opts.Get("port")
	} else {
		log.Fatal("No port specified")
		os.Exit(1)
	}

	// Get address to listen
	if opts.GetBool("listen") {
		if addr := net.ParseIP(opts.Get("listen")); addr != nil {
			listen = addr.String()
		} else {
			log.Fatal("Invalid IP address")
		}
	}

	if listen == "" {
        log.Debug("Listen on all interfaces")
		address = ":" + port
	} else {
        log.Debug("Listen on address " + listen)
		address = listen + ":" + port
	}

	if opts.GetBool("ssl") {
        log.Debug("SSL enabled")
		cfg, err := ini.Load([]byte{}, opts.Get("ssl"))
		if err != nil {
			// handle error
			log.Fatalf("Error: %v\n", err)
		}

		client.Serve(address, log, cfg)
	} else {
		client.Serve(address, log, nil)
	}
}
