package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// MailConfig defines parameters for sending mail
type MailConfig struct {
	SMTPURL  string
	SMTPPort string
	SMTPFrom string
	Username string
	Password string
}

// Config defines general config parameters
type Config struct {
	Port             string
	AllowList        []string
	LogLevel         string
	IPHeader         string
	RequireNonce     bool
	RateLimitSeconds int
	MaxAgeSeconds    int
	Secret           string
}

// GetMailConfig returns a MailConfig instance
func GetMailConfig() MailConfig {
	return MailConfig{
		SMTPURL:  getOrPanic("SMTP_URL"),
		SMTPPort: getOr("SMTP_PORT", "465"),
		SMTPFrom: getOr("SMTP_FROM", ""),
		Username: getOrPanic("SMTP_USERNAME"),
		Password: getOrPanic("SMTP_PASSWORD"),
	}
}

// GetConfig returns a Config instance
func GetConfig() Config {
	rateLimitSeconds, err := strconv.Atoi(getOr("RATE_LIMIT_SECONDS", "60"))
	if err != nil {
		log.Fatal("Could not read RATE_LIMIT_SECONDS")
		panic(err)
	}
	maxAgeSeconds, err := strconv.Atoi(getOr("MAX_AGE_SECONDS", "60"))
	if err != nil {
		log.Fatal("Could not read MAX_AGE_SECONDS")
		panic(err)
	}

	return Config{
		Port:             getOr("PORT", "8000"),
		AllowList:        strings.Split(getOrPanic("ALLOW_LIST"), ","),
		LogLevel:         getOr("LOG_LEVEL", "info"),
		IPHeader:         getOr("IP_HEADER", ""),
		RequireNonce:     getOr("REQUIRE_NONCE", "true") == "true",
		RateLimitSeconds: rateLimitSeconds,
		MaxAgeSeconds:    maxAgeSeconds,
		Secret:           getOrPanic("SECRET"),
	}
}

func getOrPanic(key string) string {
	value, ex := os.LookupEnv(key)

	if !ex {
		panic(fmt.Sprintf("Could not find variable %s", key))
	}

	return value
}

func getOr(key string, value string) string {
	value, ex := os.LookupEnv(key)

	if !ex {
		return value
	}

	return value
}
