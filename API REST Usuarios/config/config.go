package config

import (
	"os"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// SecurityConfig configuración de seguridad
type SecurityConfig struct {
	EncryptionKey string // Clave para encriptación AES
	JWTSecret     string // Secret para tokens JWT
	BCryptCost    int    // Costo para bcrypt
}

// GetBCryptCost retorna el costo de bcrypt
func (sc SecurityConfig) GetBCryptCost() int {
	return sc.BCryptCost
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			ReadTimeout:  15,
			WriteTimeout: 15,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "users_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", "32-byte-encryption-key-here!!"), // En producción debe ser segura
			JWTSecret:     getEnv("JWT_SECRET", "jwt-secret-key-change-in-production"),
			BCryptCost:    10,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
