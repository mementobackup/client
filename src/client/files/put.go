/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"bitbucket.org/ebianchi/memento-common/common"
	"github.com/op/go-logging"
	"net"
	"os"
	"time"
	"runtime"
)

func Put(logger *logging.Logger, conn net.Conn, command *common.JSONCommand) {
	var err error

	_, err = os.Stat(command.Element.Name)
	if err == nil || os.IsExist(err) {
		os.Rename(command.Element.Name, command.Element.Name+"."+time.Now().String())
	}

	switch command.Element.Type {
	case "directory":
		os.Mkdir(command.Element.Name, 0755)
	case "file":
		// TODO: download file
	case "symlink":
		// TODO: create symlink
	}

	if runtime.GOOS != "windows" {
		// TODO: chown directory
		os.Chmod(command.Element.Name, perms(command.Element.Mode))
	}
}
