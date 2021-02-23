package main

import (
	"github.com/apex/log"
	"github.com/joho/godotenv"
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
	request.HandleRequests()
}
