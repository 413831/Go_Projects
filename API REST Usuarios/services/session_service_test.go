package services

import (
	"testing"
	"time"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

func TestSessionService_CreateSession(t *testing.T) {
	sessionRepo := repositories.NewMockSessionRepositoryImpl()
	userRepo := repositories.NewMockUserRepositoryImpl()
	publisher := NewEventPublisher()

	service := NewSessionService(sessionRepo, userRepo, publisher)

	// Crear usuario primero
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Active:   true,
		Deleted:  false,
	}
	userRepo.Create(user)

	// Crear sesión
	session, err := service.CreateSession(1, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Error al crear sesión: %v", err)
	}

	if session.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", session.UserID)
	}

	if session.Token == "" {
		t.Error("Se esperaba que la sesión tuviera un token")
	}

	if !session.Active {
		t.Error("Se esperaba que la sesión estuviera activa")
	}
}

func TestSessionService_ValidateSession(t *testing.T) {
	sessionRepo := repositories.NewMockSessionRepositoryImpl()
	userRepo := repositories.NewMockUserRepositoryImpl()
	publisher := NewEventPublisher()

	service := NewSessionService(sessionRepo, userRepo, publisher)

	// Crear sesión
	session, _ := service.CreateSession(1, "127.0.0.1", "test-agent")

	// Validar sesión
	validatedSession, err := service.ValidateSession(session.Token)
	if err != nil {
		t.Fatalf("Error al validar sesión: %v", err)
	}

	if validatedSession.ID != session.ID {
		t.Errorf("Expected session ID %d, got %d", session.ID, validatedSession.ID)
	}
}

func TestSessionService_ValidateSession_Expired(t *testing.T) {
	sessionRepo := repositories.NewMockSessionRepositoryImpl()
	userRepo := repositories.NewMockUserRepositoryImpl()
	publisher := NewEventPublisher()

	service := NewSessionService(sessionRepo, userRepo, publisher)

	// Crear sesión expirada manualmente
	session := &models.Session{
		UserID:    1,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expiró hace 1 hora
		Active:    true,
	}
	sessionRepo.CreateSession(session)

	// Intentar validar sesión expirada
	_, err := service.ValidateSession("expired-token")
	if err == nil {
		t.Error("Se esperaba error al validar sesión expirada")
	}
}

func TestSessionService_GetUserSessions(t *testing.T) {
	sessionRepo := repositories.NewMockSessionRepositoryImpl()
	userRepo := repositories.NewMockUserRepositoryImpl()
	publisher := NewEventPublisher()

	service := NewSessionService(sessionRepo, userRepo, publisher)

	// Crear múltiples sesiones
	service.CreateSession(1, "127.0.0.1", "agent1")
	service.CreateSession(1, "127.0.0.2", "agent2")

	// Obtener sesiones
	sessions, err := service.GetUserSessions(1)
	if err != nil {
		t.Fatalf("Error al obtener sesiones: %v", err)
	}

	if len(sessions) < 2 {
		t.Errorf("Expected at least 2 sessions, got %d", len(sessions))
	}
}
