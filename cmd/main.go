package main

import (
	"teleport-plugin-slack-access-request/internal/app"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/logging"
)

func init() {
	logging.Init()
	config.Init()
}

func main() {
	app.Run()
}
