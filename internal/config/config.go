package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Debug    bool
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
	DSN      string
}

func LoadConfig() Config {
	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))
	dbConfig := DatabaseConfig{
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Name:     getEnv("DB_NAME", "gokku"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		DSN:      getEnv("DATABASE_DSN", getEnv("DATABASE_URL", "")),
	}

	return Config{
		Debug: debug,
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: dbConfig,
	}
}

func (dbConfig *DatabaseConfig) AssembleDSN() string {
	if dbConfig.DSN != "" {
		return dbConfig.DSN
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port, dbConfig.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
