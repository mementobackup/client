/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package client

import (
	"crypto/rand"
	"crypto/tls"
	"github.com/go-ini/ini"
	"github.com/op/go-logging"
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

func Serve(addr string, log *logging.Logger, ssl *ini.File) {
	var cmd []byte
	var ln net.Listener

	if ssl != nil {
		key := ssl.Section("ssl").Key("key").String()
		private := ssl.Section("ssl").Key("private").String()
		ln = tlsserve(addr, key, private)
        log.Debug("Opened SSL socket")
	} else {
		ln = plainserve(addr)
        log.Debug("Opened plain socket")
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			// handle error
			log.Error("Error: %v\n", err)
		}
        log.Debug("Connection from " + conn.RemoteAddr().String() + " accepted")

		cmd, err = ioutil.ReadAll(conn)
        log.Debug("Remote data readed")
		Parse(cmd, log, conn)
		conn.Close()
	}
}
