/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"bitbucket.org/ebianchi/memento-common/common"
	"encoding/hex"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

var connection net.Conn

func visitfile(fp string, fi os.FileInfo, err error) error {
	// TODO: extract ACL

	file := common.JSONFile{}

	if err != nil {
		res := common.JSONResult{"ko", err.Error()}
		res.Send(connection)
		return nil
	}

	// Set the file name and the operating system
	file.Name = fp
	file.Os = runtime.GOOS
	file.Mtime = fi.ModTime().Unix()

	if runtime.GOOS == "linux" {
		file.User = posix_user(fi)
		file.Group = posix_group(fi)
		file.Mode = fi.Mode().String()
	}

	// Set type of element (file or directory)
	if fi.IsDir() {
		file.Type = "directory"
	} else {
		file.Type = "file"
		file.Size = fi.Size()
		file.Hash = hex.EncodeToString(common.Md5(fp))
	}

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
