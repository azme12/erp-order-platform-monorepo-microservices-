package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	NATS     NATSConfig
	JWT      JWTConfig
	Services ServicesConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type NATSConfig struct {
	URL string
}

type JWTConfig struct {
	Secret              string
	UserExpirationHours int
}

type ServicesConfig struct {
	Auth      ServiceURLConfig
	Contact   ServiceURLConfig
	Inventory ServiceURLConfig
	Sales     ServiceURLConfig
	Purchase  ServiceURLConfig
}

type ServiceURLConfig struct {
	URL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", getEnv("GATEWAY_PORT", "8000")),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "microservice"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "microservice"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", ""),
			UserExpirationHours: getEnvInt("JWT_USER_EXPIRATION_HOURS", 24),
		},
		Services: ServicesConfig{
			Auth: ServiceURLConfig{
				URL: getEnv("AUTH_SERVICE_URL", "http://localhost:8002"),
			},
			Contact: ServiceURLConfig{
				URL: getEnv("CONTACT_SERVICE_URL", "http://localhost:8001"),
			},
			Inventory: ServiceURLConfig{
				URL: getEnv("INVENTORY_SERVICE_URL", "http://localhost:8003"),
			},
			Sales: ServiceURLConfig{
				URL: getEnv("SALES_SERVICE_URL", "http://localhost:8004"),
			},
			Purchase: ServiceURLConfig{
				URL: getEnv("PURCHASE_SERVICE_URL", "http://localhost:8005"),
			},
		},
	}

	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	if value := viper.GetInt(key); value != 0 {
		return value
	}
	return defaultValue
}
