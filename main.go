package main

import (
	"github.com/apex/log"
	"github.com/joho/godotenv"
	"inheaden.io/services/simple-mail-api/pkg/config"
	"inheaden.io/services/simple-mail-api/pkg/request"
)

func init() {
	setupEnv()
}

func setupEnv() {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error when reading .env file: %s", err.Error())
	}
}

func main() {
	generalConfig := config.GetConfig()
	log.SetLevelFromString(generalConfig.LogLevel)

	log.Infof("Starting service using %s log level", generalConfig.LogLevel)

	request.HandleRequests()
}
