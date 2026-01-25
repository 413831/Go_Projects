package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mypassword123"
	hash, err := HashPassword(password, 10)
	if err != nil {
		t.Fatalf("Error al hashear contraseña: %v", err)
	}

	if hash == password {
		t.Error("El hash no debería ser igual a la contraseña")
	}

	if len(hash) == 0 {
		t.Error("El hash no debería estar vacío")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mypassword123"
	hash, _ := HashPassword(password, 10)

	// Contraseña correcta
	if !CheckPassword(password, hash) {
		t.Error("La contraseña debería ser válida")
	}

	// Contraseña incorrecta
	if CheckPassword("wrongpassword", hash) {
		t.Error("La contraseña incorrecta no debería ser válida")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "mypassword123"
	hash1, _ := HashPassword(password, 10)
	hash2, _ := HashPassword(password, 10)

	// Cada hash debería ser diferente debido a la sal
	if hash1 == hash2 {
		t.Error("Los hashes deberían ser diferentes")
	}

	// Pero ambos deberían validar la misma contraseña
	if !CheckPassword(password, hash1) {
		t.Error("El primer hash debería validar la contraseña")
	}

	if !CheckPassword(password, hash2) {
		t.Error("El segundo hash debería validar la contraseña")
	}
}
