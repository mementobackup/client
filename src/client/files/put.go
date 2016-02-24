/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"github.com/mementobackup/common/src/common"
	"github.com/op/go-logging"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func fsSetAttrs(log *logging.Logger, command *common.JSONCommand) common.JSONResult {
	var result common.JSONResult

	if runtime.GOOS != "windows" {
		if res := fsPosixSetPerms(log, &command.Element); res.Result == "ko" {
			result = res
		} else {
			result = fsPosixSetAcls(log, &command.Element.Name, &command.Element.Acl)
		}
	} else {
		result = fsWindowsSetAcls(log, &command.Element.Name, &command.Element.Acl)
	}

	return result
}

func fsWindowsSetAcls(log *logging.Logger, filename *string, acls *[]common.JSONFileAcl) common.JSONResult {
	var result common.JSONResult

	// TODO: write code to set ACLs on Windows
	result = common.JSONResult{Result: "ok"}

	return result
}

func fsPosixSetAcls(log *logging.Logger, filename *string, acls *[]common.JSONFileAcl) common.JSONResult {
	var result common.JSONResult
	var err error

	for _, item := range *acls {
		fa := FileACL(*filename)
		err = fa.Set(log, item)
	}

	if err != nil {
		result = common.JSONResult{Result: "ko"}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			message := exitError.String()

			result.Message = "Status: " + strconv.Itoa(waitStatus.ExitStatus()) + " Message: " + message
		}
	} else {
		result = common.JSONResult{Result: "ok"}
	}
	return result
}

func fsPosixSetPerms(log *logging.Logger, element *common.JSONFile) common.JSONResult {
	var uid, gid int
	var result common.JSONResult

	uname, err := user.Lookup(element.User)
	if err != nil {
		uname, _ := user.Lookup("nobody")
		uid, _ = strconv.Atoi(uname.Uid)
		gid, _ = strconv.Atoi(uname.Gid)
	} else {
		uid, _ = strconv.Atoi(uname.Uid)
		gid, err = getGroupId(element.Group)
		if err != nil {
			gid, _ = getGroupId("nogroup")
		}
	}

	// TODO: fix possible error if element.Mode is empty
	if err := os.Chmod(element.Name, getPerms(element.Mode)); err != nil {
		log.Debug("Error:", err)
		result = common.JSONResult{Result: "ko", Message: err.Error()}
	} else {
		if err := os.Chown(element.Name, uid, gid); err != nil {
			log.Debug("Error:", err)
			result = common.JSONResult{Result: "ko", Message: err.Error()}
		} else {
			result = common.JSONResult{Result: "ok"}
		}
	}

	return result
}

func Put(log *logging.Logger, conn net.Conn, command *common.JSONCommand) {
	var res common.JSONResult

	if _, err := os.Stat(command.Element.Name); err == nil || os.IsExist(err) {
		os.Rename(command.Element.Name, command.Element.Name+"."+time.Now().String())
	}

	switch command.Element.Type {
	case "directory":
		os.Mkdir(command.Element.Name, 0755)
		res = fsSetAttrs(log, command)
	case "file":
		if hash, err := common.ReceiveFile(command.Element.Name, conn); hash != command.Element.Hash {
			log.Debug("Error: hash mismatch")
			res = common.JSONResult{Result: "ko", Message: "Hash mismatch"}
		} else if err != nil {
			log.Debug("Error:", err)
			res = common.JSONResult{Result: "ko", Message: "Error: " + err.Error()}
		} else {
			res = fsSetAttrs(log, command)
		}
	case "symlink":
		if err := os.Symlink(command.Element.Link, command.Element.Name); err != nil {
			log.Debug("Error:", err)
			res = common.JSONResult{Result: "ko", Message: "Error: " + err.Error()}
		} else {
			res = common.JSONResult{Result: "ok"}
		}
	}

	res.Send(conn)
}
