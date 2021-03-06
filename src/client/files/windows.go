// +build windows

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
	"os"
)

type FileACL string

func getUserName(fi os.FileInfo) (string, error) {
	// TODO: use this function for getting file username in the windows environment
	return "", nil
}

func getGroupName(fi os.FileInfo) (string, error) {
	// TODO: use this function for getting file group name in the windows environment
	return "", nil
}

func getGroupId(group string) (int, error) {
	// TODO: use this function for getting group id in the windows environment
	return -1, nil
}

func getCtime(fi os.FileInfo) int64 {
	// TODO: write code for getting windows ctime
	return 0
}

func getPerms(str string) os.FileMode {
	// TODO: use this function for getting file permissions in the windows environment
	return 0
}

func (f FileACL) List(log *logging.Logger) []common.JSONFileAcl {
	// TODO: write code for getting windows file ACLs
	return nil
}

func (f FileACL) Set(log *logging.Logger, acl common.JSONFileAcl) error {
	// TODO: write code for setting windows file ACLs
	return nil
}
