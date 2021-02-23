package config

import (
	"fmt"
	"os"
	"strings"
)

type MailConfig struct {
	SmtpURL  string
	SmtpPort string
	SmtpFrom string
	Username string
	Password string
}

type Config struct {
	Port      string
	AllowList []string
}

func GetMailConfig() MailConfig {
	return MailConfig{
		SmtpURL:  getOrPanic("SMTP_URL"),
		SmtpPort: getOrPanic("SMTP_PORT"),
		SmtpFrom: getOrPanic("SMTP_FROM"),
		Username: getOrPanic("SMTP_USERNAME"),
		Password: getOrPanic("SMTP_PASSWORD"),
	}
}

func GetConfig() Config {
	return Config{
		Port:      getOrPanic("PORT"),
		AllowList: strings.Split(getOrPanic("ALLOW_LIST"), ","),
	}
}

func getOrPanic(key string) string {
	value, ex := os.LookupEnv(key)

	if !ex {
		panic(fmt.Sprintf("Could not find variable %s", key))
	}

	return value
}
