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
	"github.com/op/go-logging"
	"net"
	"os"
)

func Parse(data []uint8, log *logging.Logger, conn net.Conn) {
	var cmd common.JSONMessage

	if err := json.Unmarshal(data, &cmd); err != nil {
		log.Debug("Error when parsing data: " + string(data))
		res := common.JSONResult{"ko", "Malformed command: " + err.Error()}
		res.Send(conn)
		return
	}
	log.Debug("Received data: " + string(data))

	switch cmd.Context {
	case "system":
		log.Debug("System context requested")
		if cmd.Command.Name == "exit" {
			log.Debug("Exit command requested")
            res := common.JSONResult{"ok", "Client closed"}
            res.Send(conn)
			os.Exit(0)
		} else if cmd.Command.Name == "exec" {
			log.Debug("Execute command requested")
			if err := common.ExecuteCMD(cmd.Command.Value); err != nil {
				log.Debug("Error when executing command: " + err.Error())
			}
            res := common.JSONResult{"ok", "Command executed"}
            res.Send(conn)
		}
	case "file":
		log.Debug("File context requested")
		if cmd.Command.Name == "list" {
			log.Debug("List command requested")
			files.List(conn, &cmd.Command)
		} else if cmd.Command.Name == "get" {
			log.Debug("Get command requested")
			common.Sendfile(cmd.Command.Filename, conn)
		} else if cmd.Command.Name == "put" {
			log.Debug("Put command requested")
			// TODO: Write code for file putting command
		} else {
			log.Debug("Invalid command requested: " + cmd.Command.Name)
			res := common.JSONResult{"ko", "Command unknown: " + cmd.Command.Name}
			res.Send(conn)
		}
	default:
		log.Debug("Invalid context requested: " + cmd.Context)
		res := common.JSONResult{"ko", "Context unknown: " + cmd.Context}
		res.Send(conn)
	}
}
