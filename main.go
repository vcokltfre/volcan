package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vcokltfre/volcan/src/api"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	godotenv.Load(".env", "/config/.env")
}

func main() {
	logrus.Info("Starting Volcan...")

	api.Start(os.Getenv("API_BIND"))
}
