package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Development(t *testing.T) {
	// Limpiar variables de entorno para el test
	os.Clearenv()
	os.Setenv("ENV", "development")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Error al cargar configuración en desarrollo: %v", err)
	}

	// En desarrollo, debería permitir valores por defecto
	if cfg.Security.EncryptionKey == "" {
		t.Error("EncryptionKey no debería estar vacía")
	}

	if cfg.Security.JWTSecret == "" {
		t.Error("JWTSecret no debería estar vacío")
	}
}

func TestLoadConfig_Production_Valid(t *testing.T) {
	// Limpiar variables de entorno
	os.Clearenv()
	os.Setenv("ENV", "production")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012") // 32 bytes exactos
	os.Setenv("JWT_SECRET", "this-is-a-very-secure-jwt-secret-key-with-more-than-32-chars")
	os.Setenv("DB_PASSWORD", "secure-production-password")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Error al cargar configuración en producción válida: %v", err)
	}

	if cfg.Security.EncryptionKey != "12345678901234567890123456789012" {
		t.Error("EncryptionKey debería usar el valor de la variable de entorno")
	}
}

func TestLoadConfig_Production_Invalid_DefaultEncryptionKey(t *testing.T) {
	// Limpiar variables de entorno
	os.Clearenv()
	os.Setenv("ENV", "production")
	// No establecer ENCRYPTION_KEY, debería usar el valor por defecto y fallar

	_, err := LoadConfig()
	if err == nil {
		t.Error("Se esperaba error al usar valor por defecto de ENCRYPTION_KEY en producción")
	}
}

func TestLoadConfig_Production_Invalid_DefaultJWTSecret(t *testing.T) {
	// Limpiar variables de entorno
	os.Clearenv()
	os.Setenv("ENV", "production")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("DB_PASSWORD", "secure-password")
	// No establecer JWT_SECRET, debería usar el valor por defecto y fallar

	_, err := LoadConfig()
	if err == nil {
		t.Error("Se esperaba error al usar valor por defecto de JWT_SECRET en producción")
	}
}

func TestLoadConfig_Production_Invalid_DefaultDBPassword(t *testing.T) {
	// Limpiar variables de entorno
	os.Clearenv()
	os.Setenv("ENV", "production")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("JWT_SECRET", "this-is-a-very-secure-jwt-secret-key-with-more-than-32-chars")
	// No establecer DB_PASSWORD, debería usar el valor por defecto y fallar

	_, err := LoadConfig()
	if err == nil {
		t.Error("Se esperaba error al usar valor por defecto de DB_PASSWORD en producción")
	}
}

func TestLoadConfig_Production_Invalid_EncryptionKeyLength(t *testing.T) {
	// Limpiar variables de entorno
	os.Clearenv()
	os.Setenv("ENV", "production")
	os.Setenv("ENCRYPTION_KEY", "short-key") // Menos de 32 bytes
	os.Setenv("JWT_SECRET", "this-is-a-very-secure-jwt-secret-key-with-more-than-32-chars")
	os.Setenv("DB_PASSWORD", "secure-password")

	_, err := LoadConfig()
	if err == nil {
		t.Error("Se esperaba error al usar ENCRYPTION_KEY con longitud incorrecta en producción")
	}
}

func TestLoadConfig_Production_Invalid_JWTSecretLength(t *testing.T) {
	// Limpiar variables de entorno
	os.Clearenv()
	os.Setenv("ENV", "production")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("JWT_SECRET", "short") // Menos de 32 caracteres
	os.Setenv("DB_PASSWORD", "secure-password")

	_, err := LoadConfig()
	if err == nil {
		t.Error("Se esperaba error al usar JWT_SECRET con longitud insuficiente en producción")
	}
}
