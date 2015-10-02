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
	"github.com/op/go-logging"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

var connection net.Conn
var acl bool
var log *logging.Logger

func visitfile(fp string, fi os.FileInfo, err error) error {
	var res common.JSONResult
	var file common.JSONFile

	if err != nil {
		res = common.JSONResult{Result: "ko", Message: err.Error()}
		res.Send(connection)
		return nil
	}

	// Set the file name and the operating system
	file.Name = fp
	file.Os = runtime.GOOS
	file.Mtime = fi.ModTime().Unix()

	if runtime.GOOS != "windows" {
		file.User, _ = username(fi)
		file.Group, _ = groupname(fi)
		file.Mode = fi.Mode().String()
		file.Ctime = ctime(fi)
	}

	// Set type of element (file or directory)
	if fi.IsDir() {
		file.Type = "directory"
		file.Size = fi.Size()
	} else if fi.Mode() & os.ModeSymlink == os.ModeSymlink {
		file.Type = "symlink"
		file.Size = fi.Size()

		link, err := os.Readlink(fp)
		if err != nil {
			log.Debug("Error when readlink for " + fp + ": " + err.Error())
		} else {
			file.Link = link
		}
	} else {
		file.Type = "file"
		file.Size = fi.Size()
		file.Hash = hex.EncodeToString(common.Md5(fp))
	}

	if acl && file.Type != "symlink" {
		if runtime.GOOS != "windows" {
			fa := FileACL(fp)
			file.Acl = fa.List(log)
		}
	}

	// Set result
	res.Result = "ok"
	res.Data = file
	res.Send(connection)
	return nil
}

func List(logger *logging.Logger, conn net.Conn, command *common.JSONCommand) {
	connection = conn
	acl = command.ACL
	log = logger

	if len(command.Paths) > 0 {
		for _, item := range command.Paths {
			// WARNING: filepath.Walk() is inefficient
			filepath.Walk(item, visitfile)
		}
	} else {
		res := common.JSONResult{Result: "ko", Message: "No directory specified"}
		res.Send(connection)
	}
}
