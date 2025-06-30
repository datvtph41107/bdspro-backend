package main

import (
	"os"

	"relay/config"
	"relay/server"
)

func main() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	config.InitConfig()

	server.NewWSServer(config.PORT).Start()
}
