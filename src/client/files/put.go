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
	"os/user"
	"runtime"
	"strconv"
	"time"
)

func setperms(element *common.JSONFile) {
	var uid, gid int

	uname, err := user.Lookup(element.User)
	if err != nil {
		uname, _ := user.Lookup("nobody")
		uid, _ = strconv.Atoi(uname.Uid)
		gid, _ = strconv.Atoi(uname.Gid)
	}
	uid, _ = strconv.Atoi(uname.Uid)
	gid, _ = getgroupid(element.Group)

	os.Chown(element.Name, uid, gid)
	os.Chmod(element.Name, getperms(element.Mode))
}

func Put(log *logging.Logger, conn net.Conn, command *common.JSONCommand) {
	var err error

	_, err = os.Stat(command.Element.Name)
	if err == nil || os.IsExist(err) {
		os.Rename(command.Element.Name, command.Element.Name+"."+time.Now().String())
	}

	switch command.Element.Type {
	case "directory":
		os.Mkdir(command.Element.Name, 0755)
	case "file":
		hash, err := common.Receivefile(command.Element.Name, conn)
		if hash != command.Element.Hash {
			log.Debug("Error: hash mismatch")
			res := common.JSONResult{Result: "ko", Message: "Hash mismatch"}
			res.Send(conn)
			return
		}

		if err != nil {
			log.Debug("Error:", err)
			res := common.JSONResult{Result: "ko", Message: "Error: " + err.Error()}
			res.Send(conn)
			return
		}
	case "symlink":
		// TODO: create symlink
	}

	if runtime.GOOS != "windows" {
		setperms(&command.Element)
	} else {
		// TODO: write code for set ACLs on Windows
	}

	res := common.JSONResult{Result: "ok"}
	res.Send(conn)
}
