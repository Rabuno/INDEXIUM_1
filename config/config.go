package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort    string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

// LoadConfig đọc biến môi trường set trong docker-compose
func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppPort:    getEnv("APP_PORT", ":8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "cms_db"),
	}
	return cfg, nil
}

// Helper để lấy DSN (Data Source Name) cho MySQL connection
func (c *Config) GetDSN() string {
	// Format: user:password@tcp(host:port)/dbname?parseTime=true
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
