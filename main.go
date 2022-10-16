package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	godotenv.Load(".env", "/config/.env")
}

func main() {
	logrus.Info("Starting Volcan...")
}
