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

/*
#include <pwd.h>
#include <grp.h>
#include <stdlib.h>

static int mygetpwuid_r(int uid, struct passwd *pwd, char *buf, size_t buflen, struct passwd **result) {
    	 return getpwuid_r(uid, pwd, buf, buflen, result);
}

static int mygetgrgid_r(int gid, struct group *grp, char *buf, size_t buflen, struct group **result) {
    	 return getgrgid_r(gid, grp, buf, buflen, result);
}
*/
import "C"

var connection net.Conn

func posix_user(fi os.FileInfo) string {
	var rv C.int
	var pwd C.struct_passwd
	var pwdres *C.struct_passwd
	var bufSize C.long
	var result string

	bufSize = 1024
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	uid := fi.Sys().(*syscall.Stat_t).Uid

	rv = C.mygetpwuid_r(C.int(uid), &pwd, (*C.char)(buf), C.size_t(bufSize), &pwdres)
	if rv != 0 {
		// Manage error
	}

	if pwdres != nil {
		result = C.GoString(pwd.pw_name)
	} else {
		// Manage error
	}

	return result
}

func posix_group(fi os.FileInfo) string {
	var rv C.int
	var grp C.struct_group
	var grpres *C.struct_group
	var bufSize C.long
	var result string

	bufSize = 1024
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	gid := fi.Sys().(*syscall.Stat_t).Gid

	rv = C.mygetgrgid_r(C.int(gid), &grp, (*C.char)(buf), C.size_t(bufSize), &grpres)
	if rv != 0 {
		// Manage error
	}

	if grpres != nil {
		result = C.GoString(grp.gr_name)
	} else {
		// Manage error
	}
	return result
}

func visitfile(fp string, fi os.FileInfo, err error) error {
	file := common.JSONFile{}

	if err != nil {
		res := common.JSONResult{"ko", err.Error()}
		res.Send(connection)
		return nil
	}

	// Set the file name and the operating system
	file.Name = fp
	file.Size = fi.Size()
	file.Os = runtime.GOOS

	if runtime.GOOS == "linux" {
		file.User = posix_user(fi)
		file.Group = posix_group(fi)
	}

	// Set type of element (file or directory)
	if fi.IsDir() {
		file.Type = "directory"
	} else {
		file.Type = "file"
	}

	// TODO: add hash, date, permission and ACL

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
