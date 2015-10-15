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

func Put(logger *logging.Logger, conn net.Conn, command *common.JSONCommand) {
	var uid, gid int
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
			// TODO: manage transfer error
		}

		if err != nil {
			// TODO: manage error
		}
	case "symlink":
		// TODO: create symlink
	}

	if runtime.GOOS != "windows" {
		uname, err := user.Lookup(command.Element.User)
		if err != nil {
			uname, _ := user.Lookup("nobody")
			uid, _ = strconv.Atoi(uname.Uid)
			gid, _ = strconv.Atoi(uname.Gid)
		}
		uid, _ = strconv.Atoi(uname.Uid)
		gid, _ = getgroupid(command.Element.Group)

		os.Chown(command.Element.Name, uid, gid)
		os.Chmod(command.Element.Name, getperms(command.Element.Mode))
	}
}
