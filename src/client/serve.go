/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package client

import (
	"code.google.com/p/goconf/conf"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

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

func Serve(addr string, ssl *conf.ConfigFile) {
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
		Parse(cmd, conn)
		conn.Close()
	}
}
