package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/redcon"
	"strconv"
	"strings"
)

var (
	addr  = ":6380"
	idGen *IdGen
)

func main() {
	initConfig()

	var initErr error
	idGen, initErr = NewIdGen()
	if initErr != nil {
		panic(initErr)
	}

	startErr := redcon.ListenAndServe(
		addr,
		handle,
		func(conn redcon.Conn) bool {
			return true
		},
		func(conn redcon.Conn, err error) {},
	)

	if startErr != nil {
		panic(startErr)
	}
}

func handle(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	case "get":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for get command")
			return
		}

		businessId, _ := strconv.Atoi(string(cmd.Args[1]))
		id := idGen.Next(uint32(businessId))

		conn.WriteBulk([]byte(strconv.FormatUint(id, 10)))
	case "parse":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for parse command")
			return
		}

		id, _ := strconv.ParseUint(string(cmd.Args[1]), 10, 64)

		result, err := json.Marshal(idGen.Parse(id))
		if err != nil {
			conn.WriteError(err.Error())
		} else {
			conn.WriteBulk(result)
		}
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		err := conn.Close()

		if err != nil {
			fmt.Println(err)
		}
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	}
}
