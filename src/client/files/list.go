/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package files

import (
	"bitbucket.org/ebianchi/memento-common/common"
	"fmt"
	"os"
	"path/filepath"
)

func visitfile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}

    // not a file
	if fi.IsDir() {
		fmt.Println(fp)
		return nil
	}

	fmt.Println(fp)
	return nil
}

func List(command *common.JSONCommand) {
	filepath.Walk(command.Directory, visitfile)
}
