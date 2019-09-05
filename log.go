package main

import (
	"log"
	"os"
	"time"
)

var logger *log.Logger

func initLog() {
	fileName := config.LogPath + time.Now().Format("20060102") + ".log"
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("file open error : %v", err)
	}

	logger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
}
