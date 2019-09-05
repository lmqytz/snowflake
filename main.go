package main

import (
	"encoding/json"
	"fmt"
	"github.com/bsm/redeo"
	"github.com/bsm/redeo/resp"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	initConfig()
	initLog()
}

var idGen *IdGen

func main() {
	netCard := config.NetCard
	workerIdString, err := getWorkId(netCard)
	if err != nil {
		logger.Fatalln(err)
	}

	workerId, _ := strconv.ParseUint(workerIdString, 10, 64)
	idGen, err = NewIdGen(uint32(workerId))
	if err != nil {
		logger.Fatalln(err)
	}

	var srv *redeo.Server
	srv = redeo.NewServer(nil)
	srv.HandleFunc("get", commandGet)
	srv.HandleFunc("parse", commandParse)
	srv.HandleFunc("ping", commandPing)

	lis, err := net.Listen("tcp", config.Listen)
	if err != nil {
		logger.Fatalln(err)
	}

	defer func() {
		if err := lis.Close(); err != nil {
			logger.Println(err)
		}
	}()

	var errChan = make(chan string)
	go func() {
		if err := srv.Serve(lis); err != nil {
			errChan <- err.Error()
		}
	}()

	ch := make(chan os.Signal, 1)
	go func() {
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	}()

	select {
	case errMsg := <-errChan:
		logger.Fatalln(errMsg)
	case <-ch:
		fmt.Println("receive signal, bye.")
	}
}

func commandGet(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 1 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		return
	}

	businessId, _ := strconv.Atoi(c.Arg(0).String())
	id := idGen.Next(uint32(businessId))
	w.AppendBulkString(strconv.FormatUint(id, 10))
}

func commandParse(w resp.ResponseWriter, c *resp.Command) {
	if c.ArgN() != 1 {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
		return
	}

	id, _ := strconv.ParseUint(c.Arg(0).String(), 10, 64)

	result, err := json.Marshal(idGen.Parse(id))
	if err != nil {
		w.AppendError(redeo.WrongNumberOfArgs(c.Name))
	} else {
		w.AppendInline(result)
	}
}

func commandPing(w resp.ResponseWriter, _ *resp.Command) {
	w.AppendInlineString("PONG")
}
