package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	EtcdAddr []string `json:"etcd_addr"`
}

var config Config

func initConfig() {
	buffer, _ := os.Open("config.json")
	err := json.NewDecoder(buffer).Decode(&config)

	fmt.Println(config.EtcdAddr)

	if err != nil {
		fmt.Println(err)
	}
}
