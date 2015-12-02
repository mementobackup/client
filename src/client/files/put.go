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

func fs_set_attrs(log *logging.Logger, command *common.JSONCommand) common.JSONResult {
	var result common.JSONResult

	if runtime.GOOS != "windows" {
		if res := fs_posix_set_perms(log, &command.Element); res.Result == "ko" {
			result = res
		} else {
			result = fs_posix_set_acls(log, &command.Element.Name, &command.Element.Acl)
		}
	} else {
		result = fs_windows_set_acls(log, &command.Element.Name, &command.Element.Acl)
	}

	return result
}

func fs_windows_set_acls(log *logging.Logger, filename *string, acls *[]common.JSONFileAcl) common.JSONResult {
	var result common.JSONResult

	// TODO: write code to set ACLs on Windows
	result = common.JSONResult{Result: "ok"}

	return result
}

func fs_posix_set_acls(log *logging.Logger, filename *string, acls *[]common.JSONFileAcl) common.JSONResult {
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

func fs_posix_set_perms(log *logging.Logger, element *common.JSONFile) common.JSONResult {
	var uid, gid int
	var result common.JSONResult

	uname, err := user.Lookup(element.User)
	if err != nil {
		uname, _ := user.Lookup("nobody")
		uid, _ = strconv.Atoi(uname.Uid)
		gid, _ = strconv.Atoi(uname.Gid)
	} else {
		uid, _ = strconv.Atoi(uname.Uid)
		gid, err = getgroupid(element.Group)
		if err != nil {
			gid, _ = getgroupid("nogroup")
		}
	}

	// TODO: fix possible error if element.Mode is empty
	if err := os.Chmod(element.Name, getperms(element.Mode)); err != nil {
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
		err = os.Symlink(command.Element.Link, command.Element.Name)
		if err != nil {
			log.Debug("Error:", err)
			res := common.JSONResult{Result: "ko", Message: "Error: " + err.Error()}
			res.Send(conn)
		} else {
			res := common.JSONResult{Result:"ok"}
			res.Send(conn)
		}
		return
	}

	res := fs_set_attrs(log, command)
	res.Send(conn)
}
