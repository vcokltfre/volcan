package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vcokltfre/volcan/src/api"
	"github.com/vcokltfre/volcan/src/bot/core"
	"github.com/vcokltfre/volcan/src/bot/modules/meta"
	"github.com/vcokltfre/volcan/src/config"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	godotenv.Load(".env", "/config/.env")

	if os.Getenv("DOCKER") == "true" {
		config.LoadConfig("/config/config.yml")
	} else {
		config.LoadConfig("config.yml")
	}
}

func main() {
	logrus.Info("Starting Volcan...")

	bot := core.Bot{
		Modules: []*core.Module{
			meta.MetaModule,
		},
	}

	err := bot.Build()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Finished building command trees.")

	go api.Start(os.Getenv("API_BIND"))

	err = bot.Start()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("API start dispatched.")

	logrus.Info("Volcan started!")

	select {}
}
