/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package main

import (
	"github.com/op/go-logging"
	"os"
)

func setlog(debug bool) *logging.Logger {
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
