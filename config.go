package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	NetCard string `json:"net_card"`
	LogPath string `json:"log_path"`
	Listen  string `json:"listen"`
}

var config Config

func initConfig() {
	buffer, _ := os.Open("config.json")
	err := json.NewDecoder(buffer).Decode(&config)

	if err != nil {
		panic(err)
	}
}
