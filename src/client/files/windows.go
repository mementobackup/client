// +build windows

/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"github.com/op/go-logging"
	"os"
)

type FileACL string

func getusername(fi os.FileInfo) (string, error) {
	// TODO: use this function for getting file username in the windows environment
	return "", nil
}

func getgroupname(fi os.FileInfo) (string, error) {
	// TODO: use this function for getting file group name in the windows environment
	return "", nil
}

func getgroupid(group string) (int, error) {
	// TODO: use this function for getting group id in the windows environment
	return -1, nil
}

func (f FileACL) List(log *logging.Logger) []common.JSONFileAcl {
	// TODO: write code for getting windows file ACLs
	return nil
}

func getctime(fi os.FileInfo) int64 {
	// TODO: write code for getting windows ctime
	return 0
}

func getperms(str string) os.FileMode {
	// TODO: use this function for getting file permissions in the windows environment
	return nil
}