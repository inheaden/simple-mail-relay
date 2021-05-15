package main

import (
	"github.com/apex/log"
	"github.com/inheaden/simple-mail-api/pkg/config"
	"github.com/inheaden/simple-mail-api/pkg/request"
	"github.com/joho/godotenv"
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
