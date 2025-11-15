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
	Secret                 string
	UserExpirationHours    int
	ServiceExpirationHours int
}

type ServicesConfig struct {
	Auth      ServiceConfig
	Contact   ServiceConfig
	Inventory ServiceConfig
	Sales     ServiceConfig
	Purchase  ServiceConfig
}

type ServiceConfig struct {
	URL    string
	Secret string
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
			Port: getEnv("PORT", "8000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "db"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			UserExpirationHours:    getEnvInt("JWT_USER_EXPIRATION_HOURS", 24),
			ServiceExpirationHours: getEnvInt("JWT_SERVICE_EXPIRATION_HOURS", 1),
		},
		Services: ServicesConfig{
			Auth: ServiceConfig{
				URL:    getEnv("AUTH_SERVICE_URL", "http://auth:8000"),
				Secret: getEnv("AUTH_SERVICE_SECRET", "auth-service-secret"),
			},
			Contact: ServiceConfig{
				URL:    getEnv("CONTACT_SERVICE_URL", "http://contact:8000"),
				Secret: getEnv("CONTACT_SERVICE_SECRET", "contact-service-secret"),
			},
			Inventory: ServiceConfig{
				URL:    getEnv("INVENTORY_SERVICE_URL", "http://inventory:8000"),
				Secret: getEnv("INVENTORY_SERVICE_SECRET", "inventory-service-secret"),
			},
			Sales: ServiceConfig{
				URL:    getEnv("SALES_SERVICE_URL", "http://sales:8000"),
				Secret: getEnv("SALES_SERVICE_SECRET", "sales-service-secret"),
			},
			Purchase: ServiceConfig{
				URL:    getEnv("PURCHASE_SERVICE_URL", "http://purchase:8000"),
				Secret: getEnv("PURCHASE_SERVICE_SECRET", "purchase-service-secret"),
			},
		},
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
