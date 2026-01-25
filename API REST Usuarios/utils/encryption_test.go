package utils

import (
	"testing"
)

func TestEncryptionService_EncryptDecrypt(t *testing.T) {
	service, err := NewEncryptionService("12345678901234567890123456789012") // 32 bytes exactos
	if err != nil {
		t.Fatalf("Error al crear servicio de encriptación: %v", err)
	}

	plaintext := "Información sensible"
	encrypted, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Error al encriptar: %v", err)
	}

	if encrypted == plaintext {
		t.Error("El texto encriptado no debería ser igual al texto plano")
	}

	decrypted, err := service.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Error al desencriptar: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestEncryptionService_InvalidKey(t *testing.T) {
	_, err := NewEncryptionService("short-key")
	if err == nil {
		t.Error("Se esperaba error con clave inválida")
	}
}

func TestEncryptionService_DifferentEncryptions(t *testing.T) {
	service, err := NewEncryptionService("12345678901234567890123456789012") // 32 bytes exactos
	if err != nil {
		t.Fatalf("Error al crear servicio de encriptación: %v", err)
	}

	plaintext := "test"
	encrypted1, _ := service.Encrypt(plaintext)
	encrypted2, _ := service.Encrypt(plaintext)

	// Cada encriptación debería ser diferente debido al nonce
	if encrypted1 == encrypted2 {
		t.Error("Las encriptaciones deberían ser diferentes")
	}

	// Pero ambas deberían desencriptarse al mismo valor
	decrypted1, _ := service.Decrypt(encrypted1)
	decrypted2, _ := service.Decrypt(encrypted2)

	if decrypted1 != decrypted2 {
		t.Error("Ambas encriptaciones deberían desencriptarse al mismo valor")
	}
}
