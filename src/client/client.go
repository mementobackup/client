/*
 Copyright (C) 2015 Enrico Bianchi (enrico.bianchi@gmail.com)
 Project       Memento
 Description   A backup system
 License       GPL version 2 (see GPL.txt for details)
*/

package client

import (
	"bitbucket.org/ebianchi/memento-common/common"
	"client/files"
	"encoding/json"
	"net"
	"os"
)

func Parse(data []uint8, conn net.Conn) {
	var cmd common.JSONMessage

	if err := json.Unmarshal(data, &cmd); err != nil {
		common.Sendresult(conn, "ko", "Malformed command")
		return
	}

	switch cmd.Context {
	case "system":
		if cmd.Command.Name == "exit" {
			os.Exit(0)
		} else if cmd.Command.Name == "exec" {
			// TODO: Write code for executing commands
		}
	case "file":
		if cmd.Command.Name == "list" {
			files.List(&cmd.Command)
		} else if cmd.Command.Name == "get" {
			// TODO: Write fode for file getting command
		} else if cmd.Command.Name == "put" {
			// TODO: Write fode for file putting command
		} else {
			common.Sendresult(conn, "ko", "Malformed command")
		}
	default:
		common.Sendresult(conn, "ko", "Malformed command")
	}

}
