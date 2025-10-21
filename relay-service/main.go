package main

import (
	"os"

	"relay/config"
	"relay/server"
)

func main() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	config.LoadConfig()

	server.NewWSServer(config.AppProperties.Server.Port).Start()
}
