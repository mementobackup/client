/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package client

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/go-ini/ini"
	"github.com/op/go-logging"
	"io/ioutil"
	"net"
	"time"
)

func tlsserve(log *logging.Logger, addr, certificate, key string) net.Listener {
	var err error
	var ln net.Listener

	log.Debug("SSL certificate: %s", certificate)
	log.Debug("SSL private key: %s", key)

	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		log.Fatalf("Error when instantiate SSL connection: %v", err)
	}

	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := tls.Config{
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequestClientCert,
	}

	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader

	ln, err = tls.Listen("tcp", addr, &config)
	if err != nil {
		log.Fatalf("Error when open connecting with host: %v", err)
	}

	return ln
}

func plainserve(log *logging.Logger, addr string) net.Listener {
	var err error
	var ln net.Listener

	ln, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return ln
}

func Serve(log *logging.Logger, addr string, ssl *ini.File) {
	var cmd []byte
	var ln net.Listener

	if ssl != nil {
		certificate := ssl.Section("ssl").Key("certificate").String()
		key := ssl.Section("ssl").Key("key").String()

		ln = tlsserve(log, addr, certificate, key)
		log.Debug("Opened SSL socket")
	} else {
		ln = plainserve(log, addr)
		log.Debug("Opened plain socket")
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		buff := bufio.NewReader(conn)

		if err != nil {
			// handle error
			log.Error("Error: %v\n", err)
		}
		log.Debug("Connection from " + conn.RemoteAddr().String() + " accepted")

		cmd, err = buff.ReadBytes('\n')
		log.Debug("Remote data readed")
		Parse(log, cmd, conn)
		conn.Close()
	}
}
