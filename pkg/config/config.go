package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port                int
	GracefulShutdownSec int
}

type AppConfig struct {
	DB       DBConfig
	Server   ServerConfig
	Auth     AuthConfig
	Polar    PolarConfig
	LogLevel string
	Env      string
}

type AuthConfig struct {
	JwksURL string
}

type PolarConfig struct {
	AccessToken string
}

func LoadConfig() (*AppConfig, error) {
	portStr := os.Getenv("DB_PORT")
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		portInt = 5432
	}
	config := &AppConfig{
		DB: DBConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     portInt,
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_DATABASE"),
		},
		Server: ServerConfig{
			Port:                8080,
			GracefulShutdownSec: 5,
		},
		Auth: AuthConfig{
			JwksURL: os.Getenv("JWKS_URL"),
		},
		Polar: PolarConfig{
			AccessToken: os.Getenv("POLAR_ACCESS_TOKEN"),
		},
		LogLevel: "info",
		Env:      os.Getenv("APP_ENV"),
	}
	return config, nil
}
