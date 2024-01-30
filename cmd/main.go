package main

import (
	commontConfig "github.com/quadev-ltd/qd-common/pkg/config"

	"qd-email-api/internal/application"
	"qd-email-api/internal/config"
)

func main() {

	var configurations config.Config
	configLocation := "./internal/config"
	configurations.Load(configLocation)

	var centralConfig commontConfig.Config
	centralConfig.Load(
		configurations.Environment,
		configurations.AWS.Key,
		configurations.AWS.Secret,
	)

	application := application.NewApplication(&configurations, &centralConfig)
	application.StartServer()

	defer application.Close()
}
