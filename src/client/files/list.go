/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"encoding/hex"
	"github.com/mementobackup/common/src/common"
	"github.com/op/go-logging"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type visit struct {
	connection net.Conn
	acl        bool
	exclude    string
	log        *logging.Logger
}

// Ugly but necessary because filepath.Match stop it when encounter filepath.Separator
func ismatch(pattern, item string) bool {
	elements := strings.Split(item, string(filepath.Separator))
	for _, element := range elements {
		matched, _ := filepath.Match(pattern, element)
		if matched {
			return true
		}
	}
	return false
}

func (v visit) visitfile(fp string, fi os.FileInfo, err error) error {
	var res common.JSONResult
	var file common.JSONFile

	if err != nil {
		res = common.JSONResult{Result: "ko", Message: err.Error()}
		res.Send(v.connection)
		return nil
	}

	if ismatch(v.exclude, fp) {
		return nil
	}

	// Set the file name and the operating system
	file.Name = fp
	file.Os = runtime.GOOS
	file.Mtime = fi.ModTime().Unix()

	if runtime.GOOS != "windows" {
		file.User, _ = getusername(fi)
		file.Group, _ = getgroupname(fi)
		file.Mode = fi.Mode().String()
		file.Ctime = getctime(fi)
	}

	// Set type of element (file or directory)
	if fi.IsDir() {
		file.Type = "directory"
		file.Size = fi.Size()
	} else if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		file.Type = "symlink"
		file.Size = fi.Size()

		link, err := os.Readlink(fp)
		if err != nil {
			v.log.Debug("Error when readlink for " + fp + ": " + err.Error())
		} else {
			file.Link = link
		}
	} else {
		file.Type = "file"
		file.Size = fi.Size()
		file.Hash = hex.EncodeToString(common.Md5(fp))
	}

	if v.acl && file.Type != "symlink" {
		if runtime.GOOS != "windows" {
			fa := FileACL(fp)
			file.Acl = fa.List(v.log)
		}
	}

	// Set result
	res.Result = "ok"
	res.Data = file
	res.Send(v.connection)
	return nil
}

func List(logger *logging.Logger, conn net.Conn, command *common.JSONCommand) {
	v := visit{
		connection: conn,
		acl:        command.ACL,
		exclude:    command.Exclude,
		log:        logger,
	}

	if len(command.Paths) > 0 {
		for _, path := range command.Paths {
			filepath.Walk(path, v.visitfile)
		}
	} else {
		res := common.JSONResult{Result: "ko", Message: "No directory specified"}
		res.Send(v.connection)
	}
}
