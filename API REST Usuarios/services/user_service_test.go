package services

import (
	"testing"
	"time"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

func TestUserService_Create(t *testing.T) {
	userRepo := repositories.NewMockUserRepositoryImpl()
	roleRepo := repositories.NewMockRoleRepositoryImpl()
	piiRepo := repositories.NewMockPIIRepositoryImpl()

	hasher := NewBCryptPasswordHasher(10)
	encryptor, _ := NewAESEncryptor("12345678901234567890123456789012")
	validator := NewUserValidator()
	publisher := NewEventPublisher()

	service := NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, publisher)

	req := &models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, err := service.Create(req)
	if err != nil {
		t.Fatalf("Error al crear usuario: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	if user.Password == "password123" {
		t.Error("Password no debería estar en texto plano")
	}
}

func TestUserService_GetByID(t *testing.T) {
	userRepo := repositories.NewMockUserRepositoryImpl()
	roleRepo := repositories.NewMockRoleRepositoryImpl()
	piiRepo := repositories.NewMockPIIRepositoryImpl()

	hasher := NewBCryptPasswordHasher(10)
	encryptor, _ := NewAESEncryptor("12345678901234567890123456789012")
	validator := NewUserValidator()
	publisher := NewEventPublisher()

	service := NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, publisher)

	// Crear usuario primero
	req := &models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	// Obtener usuario
	user, err := service.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Error al obtener usuario: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("Expected ID %d, got %d", createdUser.ID, user.ID)
	}
}

func TestUserService_Update(t *testing.T) {
	userRepo := repositories.NewMockUserRepositoryImpl()
	roleRepo := repositories.NewMockRoleRepositoryImpl()
	piiRepo := repositories.NewMockPIIRepositoryImpl()

	hasher := NewBCryptPasswordHasher(10)
	encryptor, _ := NewAESEncryptor("12345678901234567890123456789012")
	validator := NewUserValidator()
	publisher := NewEventPublisher()

	service := NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, publisher)

	// Crear usuario
	req := &models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	// Actualizar usuario
	newEmail := "newemail@example.com"
	updateReq := &models.UpdateUserRequest{
		Email: &newEmail,
	}

	updatedUser, err := service.Update(createdUser.ID, updateReq)
	if err != nil {
		t.Fatalf("Error al actualizar usuario: %v", err)
	}

	if updatedUser.Email != newEmail {
		t.Errorf("Expected email '%s', got '%s'", newEmail, updatedUser.Email)
	}
}

func TestUserService_Delete(t *testing.T) {
	userRepo := repositories.NewMockUserRepositoryImpl()
	roleRepo := repositories.NewMockRoleRepositoryImpl()
	piiRepo := repositories.NewMockPIIRepositoryImpl()

	hasher := NewBCryptPasswordHasher(10)
	encryptor, _ := NewAESEncryptor("12345678901234567890123456789012")
	validator := NewUserValidator()
	publisher := NewEventPublisher()

	service := NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, publisher)

	// Crear usuario
	req := &models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	// Eliminar usuario
	err := service.Delete(createdUser.ID)
	if err != nil {
		t.Fatalf("Error al eliminar usuario: %v", err)
	}

	// Intentar obtener usuario eliminado
	_, err = service.GetByID(createdUser.ID)
	if err == nil {
		t.Error("Se esperaba error al obtener usuario eliminado")
	}
}

func TestUserService_GrantRole(t *testing.T) {
	userRepo := repositories.NewMockUserRepositoryImpl()
	roleRepo := repositories.NewMockRoleRepositoryImpl()
	piiRepo := repositories.NewMockPIIRepositoryImpl()

	hasher := NewBCryptPasswordHasher(10)
	encryptor, _ := NewAESEncryptor("12345678901234567890123456789012")
	validator := NewUserValidator()
	publisher := NewEventPublisher()

	service := NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, publisher)

	// Crear usuario
	req := &models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	// Otorgar rol
	err := service.GrantRole(createdUser.ID, 1, createdUser.ID)
	if err != nil {
		t.Fatalf("Error al otorgar rol: %v", err)
	}

	// Verificar que el rol fue otorgado
	user, _ := service.GetByID(createdUser.ID)
	if len(user.Roles) == 0 {
		t.Error("Se esperaba que el usuario tuviera al menos un rol")
	}
}

func TestUserService_SavePII(t *testing.T) {
	userRepo := repositories.NewMockUserRepositoryImpl()
	roleRepo := repositories.NewMockRoleRepositoryImpl()
	piiRepo := repositories.NewMockPIIRepositoryImpl()

	hasher := NewBCryptPasswordHasher(10)
	encryptor, _ := NewAESEncryptor("12345678901234567890123456789012")
	validator := NewUserValidator()
	publisher := NewEventPublisher()

	service := NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, publisher)

	// Crear usuario
	req := &models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	createdUser, _ := service.Create(req)

	// Guardar PII
	pii := &models.PII{
		UserID:      createdUser.ID,
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		SSN:         "123-45-6789",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	err := service.SavePII(pii)
	if err != nil {
		t.Fatalf("Error al guardar PII: %v", err)
	}

	// Verificar que el PII fue guardado y encriptado
	user, _ := service.GetByID(createdUser.ID)
	if user.PII == nil {
		t.Error("Se esperaba que el usuario tuviera PII")
	}

	if user.PII.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got '%s'", user.PII.FirstName)
	}
}
