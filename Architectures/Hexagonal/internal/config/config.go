package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	ExternalAPIURL string
}

func Load() (*Config, error) {
	// Cargar variables de entorno desde .env si existe
	_ = godotenv.Load()

	config := &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		ExternalAPIURL: getEnv("EXTERNAL_API_URL", "http://localhost:3000/api"),
	}

	if config.ExternalAPIURL == "" {
		return nil, fmt.Errorf("EXTERNAL_API_URL es requerida")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
