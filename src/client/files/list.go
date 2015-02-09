/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"bitbucket.org/ebianchi/memento-common/common"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

var connection net.Conn

func visitfile(fp string, fi os.FileInfo, err error) error {
	file := common.JSONFile{}

	if err != nil {
		res := common.JSONResult{"ko", err.Error()}
		res.Send(connection)
		return nil
	}

	// Set the file name and the operating system
	file.Name = fp
	file.Os = runtime.GOOS

	if runtime.GOOS == "linux" {
		// FIXME: Convert UID into Username (you need to use a C call)
		//file.User = fi.Sys().(*syscall.Stat_t).Uid

		// FIXME: Convert GID into Groupname (you need to use a C call)
		//file.Group = fi.Sys().(*syscall.Stat_t).Gid
	}

	// Set type of element (file or directory)
	if fi.IsDir() {
		file.Type = "directory"
	} else {
		file.Type = "file"
	}

	// TODO: add hash, size, date, permission and ACL

	// Set result
	file.Result = "ok"
	file.Send(connection)
	return nil
}

func List(conn net.Conn, command *common.JSONCommand) {
	connection = conn

	if command.Directory != "" {
		// WARNING: filepath.Walk() is inefficient
		filepath.Walk(command.Directory, visitfile)
	} else {
		res := common.JSONResult{"ko", "No directory specified"}
		res.Send(connection)
	}
}
