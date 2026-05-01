package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	JWTRefreshSecret  string
}

func LoadConfig() *Config {
	return &Config{
		AppPort:          getEnv("APP_PORT", "8083"),
		DBHost:           getEnv("DB_HOST", "192.168.1.18"),
		DBPort:           getEnv("DB_PORT", "5433"),
		DBUser:           getEnv("DB_USER", ""),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", "zoztool"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
	}
}

// Validate panics at startup if required secrets are missing.
func (c *Config) Validate() {
	if c.JWTSecret == "" {
		panic("JWT_SECRET is required but not set")
	}
	if c.JWTRefreshSecret == "" {
		panic("JWT_REFRESH_SECRET is required but not set")
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
