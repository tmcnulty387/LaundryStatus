package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL      string
	Port             int
	Host             string
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string
	SMSEnabled       bool
}

func Load() *Config {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		port = 8080
	}

	return &Config{
		DatabaseURL:      getEnv("GOOSE_DBSTRING", ""),
		Host:             getEnv("HOST", "localhost"),
		Port:             port,
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
		SMSEnabled:       getEnv("SMS_ENABLED", "false") == "true",
	}
}

func (c *Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable %s not found, using fallback", key)
	return fallback
}
