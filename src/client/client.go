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
		res := common.JSONResult{"ko", "Malformed command: " + err.Error()}
		res.Send(conn)
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
			files.List(conn, &cmd.Command)
		} else if cmd.Command.Name == "get" {
			common.Sendfile(cmd.Command.Filename, conn)
		} else if cmd.Command.Name == "put" {
			// TODO: Write fode for file putting command
		} else {
			res := common.JSONResult{"ko", "Command unknown: " + cmd.Command.Name}
			res.Send(conn)
		}
	default:
		res := common.JSONResult{"ko", "Context unknown: " + cmd.Context}
		res.Send(conn)
	}

}
