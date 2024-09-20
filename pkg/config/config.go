package config

import (
	"os"
)

type Config struct {
	JWTSecret        string
	DBUser           string
	DBPass           string
	DBName           string
	CurrentIPAddress string
	EmailUser        string
	EmailPassword    string
}

func NewConfig() *Config {
	return &Config{
		JWTSecret:        getEnv("JWT_SECRET", "defaultsecret"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPass:           getEnv("DB_PASS", "1915"),
		DBName:           getEnv("DB_NAME", "authservice"),
		CurrentIPAddress: getEnv("CURRENT_IP", "127.0.0.1"),
		EmailUser:        getEnv("EMAIL_USER", "besthacker8163264@gmail.com"),
		EmailPassword:    getEnv("EMAIL_PASSWORD", "8163264128"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
