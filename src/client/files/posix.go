// +build !windows

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
	"github.com/op/go-logging"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type FileACL string

func username(fi os.FileInfo) (string, error) {
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

func groupname(fi os.FileInfo) (string, error) {
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

func (f FileACL) List(log *logging.Logger) []common.JSONFileAcl {
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
				acl.Mode = strings.Split(line, ":")[2][:3]
				result = append(result, acl)
			}
			if line[:5] == "group" && line[5:7] != "::" {
				acl.Group = strings.Split(line, ":")[1]
				acl.Mode = strings.Split(line, ":")[2][:3]
				result = append(result, acl)
			}
		}
	}

	return result
}

func ctime(fi os.FileInfo) int64 {
	var result int64

	stat := fi.Sys().(*syscall.Stat_t)
	result = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)).Unix()

	return result
}

func convert(perms string) os.FileMode {
	var others, group, user int
	var sticky, sgid, suid bool
	var result os.FileMode

	others = 0
	group = 0
	user = 0

	sticky = false
	suid = false
	sgid = false

	compute := func(perm string) int {
		if perm == "r" {
			return 4
		} else if perm == "w" {
			return 2
		} else if perm == "x" || perm == "s" || perm == "t" {
			return 1
		} else {
			return 0
		}
	}

	for i := 0; i < len(perms); i++ {
		perm := string(perms[len(perms)-1-i])
		if perm != "-" {
			if i == 0 || i == 1 || i == 2 {
				if perm == "t" {
					others += compute(perm)
					sticky = true
				} else if perm == "T" {
					sticky = true
				} else {
					others += compute(perm)
				}
			}

			if i == 3 || i == 4 || i == 5 {
				if perm == "s" {
					group += compute(perm)
					sgid = true
				} else if perm == "S" {
					sgid = true
				} else {
					group += compute(perm)
				}
			}

			if i == 6 || i == 7 || i == 8 {
				if perm == "s" {
					user += compute(perm)
					suid = true
				} else if perm == "S" {
					suid = true
				} else {
					user += compute(perm)
				}

			}
		}
	}
	octal, _ := strconv.ParseInt(strconv.Itoa(user)+strconv.Itoa(group)+strconv.Itoa(others), 8, 32)

	result = os.FileMode(octal)

	if suid {
		result = result | os.ModeSetuid
	}

	if sgid {
		result = result | os.ModeSetgid
	}

	if sticky {
		result = result | os.ModeSticky
	}

	return result
}
