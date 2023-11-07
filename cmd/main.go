package main

import (
	"qd_email_api/internal/application"
	"qd_email_api/internal/config"
)

func main() {

	var configurations config.Config
	configLocation := "./internal/config"
	configurations.Load(configLocation)

	application := application.NewApplication(&configurations)
	application.StartServer()

	defer application.Close()
}
