// +build windows

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
	"os"
)

type FileACL string

func username(fi os.FileInfo) (string, error) {
	// TODO: use this function for getting file username in the windows environment
	return "", nil
}

func groupname(fi os.FileInfo) (string, error) {
	// TODO: use this function for getting file username in the windows environment
	return "", nil
}

func (f FileACL) List(log *logging.Logger) []common.JSONFileAcl {
	// TODO: write code for getting windows file ACLs
	return nil
}

func ctime(fi os.FileInfo) int64 {
	// TODO: write code for getting windows ctime
	return 0
}