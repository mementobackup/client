/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

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

import (
	"bitbucket.org/ebianchi/memento-common/common"
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func posix_user(fi os.FileInfo) (string, error) {
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
		return "", errors.New("Could not read username")
	}

	if pwdres != nil {
		result = C.GoString(pwd.pw_name)
	} else {
		return "", errors.New("Could not convert username")
	}

	return result, nil
}

func posix_group(fi os.FileInfo) (string, error) {
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
		return "", errors.New("Could not read groupname")
	}

	if grpres != nil {
		result = C.GoString(grp.gr_name)
	} else {
		return "", errors.New("Could not convert groupname")
	}
	return result, nil
}

type FileACL string

func (f FileACL) List() []common.JSONFileAcl {
	var result []common.JSONFileAcl
	var out bytes.Buffer

	process := exec.Command("getfacl", string(f))

	process.Stdout = &out
	err := process.Run()
	if err != nil {
		log.Fatal(err)
	}
	output := strings.Split(out.String(), "\n")

	for _, line := range output {
		var acl common.JSONFileAcl
		if len(line) > 0 {
			if line[:4] == "user" && line[4:6] != "::" {
				acl.User = strings.Split(line, ":")[1]
				acl.Mode = strings.Split(line, ":")[2]
				result = append(result, acl)
			}
			if line[:5] == "group" && line[5:7] != "::" {
				acl.Group = strings.Split(line, ":")[1]
				acl.Mode = strings.Split(line, ":")[2]
				result = append(result, acl)
			}
		}
	}

	return result
}
