package config

import "os"

// SMTP configuration variables
var (
	CONFIG_SMTP_SERVICES = getEnv("SMTP_SERVICES", "gmail")
	CONFIG_SMTP_HOST     = getEnv("SMTP_HOST", "smtp.gmail.com")
	CONFIG_SMTP_PORT     = getEnv("SMTP_PORT", "465")
	CONFIG_SMTP_USER     = getEnv("SMTP_USER", "cretivesoft@gmail.com")
	CONFIG_SMTP_PASS     = getEnv("SMTP_PASS", "ctuvemflcjsgvcpz")
)

type Email struct {
	To, Cc        []string
	Subject, Body string
}

// getEnv gets an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
