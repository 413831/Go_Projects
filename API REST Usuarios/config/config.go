package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

// Valores por defecto inseguros que no deben usarse en producción
const (
	defaultEncryptionKey = "32-byte-encryption-key-here!!"
	defaultJWTSecret     = "jwt-secret-key-change-in-production"
	defaultDBPassword    = "postgres"
)

// LoadConfig carga la configuración desde variables de entorno
// En producción, las claves de seguridad deben estar definidas como variables de entorno
func LoadConfig() (*Config, error) {
	env := getEnv("ENV", "development")
	isProduction := env == "production"

	cfg := &Config{
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
			Password: getEnv("DB_PASSWORD", defaultDBPassword),
			DBName:   getEnv("DB_NAME", "users_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", defaultEncryptionKey),
			JWTSecret:     getEnv("JWT_SECRET", defaultJWTSecret),
			BCryptCost:    getEnvAsInt("BCRYPT_COST", 10),
		},
	}

	// Validar que en producción no se usen valores por defecto inseguros
	if isProduction {
		if err := validateProductionConfig(cfg); err != nil {
			return nil, fmt.Errorf("configuración de producción inválida: %w", err)
		}
	}

	return cfg, nil
}

// validateProductionConfig valida que en producción no se usen valores por defecto inseguros
func validateProductionConfig(cfg *Config) error {
	if cfg.Security.EncryptionKey == defaultEncryptionKey {
		return errors.New("ENCRYPTION_KEY debe estar definida como variable de entorno en producción")
	}
	if len(cfg.Security.EncryptionKey) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY debe tener exactamente 32 bytes, tiene %d", len(cfg.Security.EncryptionKey))
	}

	if cfg.Security.JWTSecret == defaultJWTSecret {
		return errors.New("JWT_SECRET debe estar definida como variable de entorno en producción")
	}
	if len(cfg.Security.JWTSecret) < 32 {
		return errors.New("JWT_SECRET debe tener al menos 32 caracteres en producción")
	}

	if cfg.Database.Password == defaultDBPassword {
		return errors.New("DB_PASSWORD debe estar definida como variable de entorno en producción")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
