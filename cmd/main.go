package main

import (
	"qd-email-api/internal/application"
	"qd-email-api/internal/config"
)

func main() {

	var configurations config.Config
	configLocation := "./internal/config"
	configurations.Load(configLocation)

	application := application.NewApplication(&configurations)
	application.StartServer()

	defer application.Close()
}
