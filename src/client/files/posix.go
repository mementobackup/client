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


static int mygetgrnam_r(const char *name, struct group *grp, char *buf, size_t buflen, struct group **result) {
                 return getgrnam_r(name, grp, buf, buflen, result);
}
*/
import "C"

import (
	"bytes"
	"errors"
	"github.com/mementobackup/common/src/common"
	"github.com/op/go-logging"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type FileACL string

func getusername(fi os.FileInfo) (string, error) {
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

func getgroupname(fi os.FileInfo) (string, error) {
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

func getgroupid(group string) (int, error) {
	var rv C.int
	var grp C.struct_group
	var grpres *C.struct_group
	var bufSize C.long
	var result int

	bufSize = 1024
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	rv = C.mygetgrnam_r(C.CString(group), &grp, (*C.char)(buf), C.size_t(bufSize), &grpres)
	if rv != 0 {
		return -1, errors.New("Could not read groupid")
	}

	if grpres != nil {
		result = int(C.int(grp.gr_gid))
	} else {
		return -1, errors.New("Could not convert groupid")
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
		log.Fatal(err) // FIXME: manage error
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

func (f FileACL) Set(log *logging.Logger, acl common.JSONFileAcl) error {
	command := "setfacl"
	args := []string{}

	args = append(args, "-m")
	if acl.User != "" {
		args = append(args, "u:"+acl.User+":"+acl.Mode)
	} else {
		args = append(args, "g:"+acl.Group+":"+acl.Mode)
	}
	args = append(args, string(f))

	process := exec.Command(command, args...)

	err := process.Run()
	if err != nil {
		log.Fatal(err) // FIXME: manage error
		return err
	} else {
		return nil
	}
}

func getctime(fi os.FileInfo) int64 {
	var result int64

	stat := fi.Sys().(*syscall.Stat_t)
	result = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)).Unix()

	return result
}

func getperms(str string) os.FileMode {
	computeperms := func(perms string) string {
		var oct int
		if string(perms[0]) == "r" {
			oct += 4
		}
		if string(perms[1]) == "w" {
			oct += 2
		}
		if string(perms[2]) == "x" {
			oct += 1
		}

		return strconv.Itoa(oct)
	}

	computemodes := func(perms int64, modes string) os.FileMode {
		result := os.FileMode(perms)
		for _, mode := range modes {
			switch string(mode) {
			case "u":
				result |= os.ModeSetuid
			case "g":
				result |= os.ModeSetgid
			case "t":
				result |= os.ModeSticky
			}
		}
		return result
	}

	mode, perms := str[:len(str)-9], str[len(str)-9:]
	user, group, others := perms[0:3], perms[3:6], perms[6:9]

	octal := computeperms(user) + computeperms(group) + computeperms(others)
	conv, _ := strconv.ParseInt(octal, 8, 32)

	result := computemodes(conv, mode)

	return result
}
