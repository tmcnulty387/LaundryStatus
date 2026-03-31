package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL       string
	GooseMigrationDir string
	GooseDriver       string
	RunMigrations     bool
	Host              string
	Port              int
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioFromNumber  string
	SMSEnabled        bool
}

func Load() *Config {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		port = 8080
	}

	sms := getEnv("SMS_ENABLED", "false") == "true"
	if sms {
		log.Println("SMS notifications are enabled")
	} else {
		log.Println("SMS notifications are disabled")
	}

	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		GooseMigrationDir: getEnv("GOOSE_MIGRATION_DIR", "/migrations"),
		GooseDriver:       getEnv("GOOSE_DRIVER", "postgres"),
		RunMigrations:     getEnv("RUN_MIGRATIONS", "false") == "true",
		Host:              getEnv("HOST", "0.0.0.0"),
		Port:              port,
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber:  getEnv("TWILIO_FROM_NUMBER", ""),
		SMSEnabled:        sms,
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
