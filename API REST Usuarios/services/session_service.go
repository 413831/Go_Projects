package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

// SessionServiceInterface define la interfaz para el servicio de sesiones
type SessionServiceInterface interface {
	CreateSession(userID int64, ipAddress, userAgent string) (*models.Session, error)
	ValidateSession(token string) (*models.Session, error)
	GetUserSessions(userID int64) ([]*models.Session, error)
	DeactivateSession(sessionID int64) error
	Stop()
}

// SessionService implementa SessionServiceInterface con patrones de diseño
type SessionService struct {
	sessionRepo    repositories.SessionRepository
	userRepo       repositories.UserRepository
	eventPublisher EventPublisher
	cleanupChannel chan bool
	cleanupWG      sync.WaitGroup
	mu             sync.Mutex
}

// NewSessionService crea una nueva instancia del servicio de sesiones
func NewSessionService(
	sessionRepo repositories.SessionRepository,
	userRepo repositories.UserRepository,
	eventPublisher EventPublisher,
) SessionServiceInterface {
	ss := &SessionService{
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
		cleanupChannel: make(chan bool),
	}
	
	// Iniciar goroutine para limpieza periódica de sesiones
	go ss.startCleanupRoutine()
	
	return ss
}

// NewSessionServiceV2 mantiene compatibilidad con application_v2/factory.
func NewSessionServiceV2(
	sessionRepo repositories.SessionRepository,
	userRepo repositories.UserRepository,
	eventPublisher EventPublisher,
) SessionServiceInterface {
	return NewSessionService(sessionRepo, userRepo, eventPublisher)
}

func (ss *SessionService) CreateSession(userID int64, ipAddress, userAgent string) (*models.Session, error) {
	token, err := generateToken()
	if err != nil {
		ss.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			UserID:  userID,
			Message: fmt.Sprintf("Error generando token: %v", err),
		}))
		return nil, err
	}

	session := &models.Session{
		UserID:    userID,
		Token:     token,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Sesión válida por 24 horas
		Active:    true,
	}

	err = ss.sessionRepo.CreateSession(session)
	if err != nil {
		ss.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			UserID:  userID,
			Message: fmt.Sprintf("Error creando sesión: %v", err),
		}))
		return nil, err
	}

	// Publicar evento de sesión creada
	ss.eventPublisher.Publish(NewUserEvent("session.created", &SessionData{
		SessionID: session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		IPAddress: session.IPAddress,
	}))

	return session, nil
}

func (ss *SessionService) ValidateSession(token string) (*models.Session, error) {
	session, err := ss.sessionRepo.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		ss.sessionRepo.DeactivateSession(session.ID)
		return nil, ErrSessionExpired
	}

	// Actualizar última actividad en goroutine para no bloquear
	go func() {
		ss.sessionRepo.UpdateSessionActivity(session.ID)
	}()

	return session, nil
}

func (ss *SessionService) GetUserSessions(userID int64) ([]*models.Session, error) {
	return ss.sessionRepo.GetUserSessions(userID)
}

func (ss *SessionService) DeactivateSession(sessionID int64) error {
	err := ss.sessionRepo.DeactivateSession(sessionID)
	if err != nil {
		ss.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			Message: fmt.Sprintf("Error desactivando sesión %d: %v", sessionID, err),
		}))
		return err
	}

	return nil
}

// startCleanupRoutine inicia una goroutine que limpia sesiones expiradas periódicamente
func (ss *SessionService) startCleanupRoutine() {
	ss.cleanupWG.Add(1)
	ticker := time.NewTicker(1 * time.Hour) // Limpiar cada hora
	defer ticker.Stop()
	defer ss.cleanupWG.Done()

	for {
		select {
		case <-ticker.C:
			ss.cleanupExpiredSessions()
		case <-ss.cleanupChannel:
			return
		}
	}
}

// cleanupExpiredSessions limpia las sesiones expiradas usando goroutines
func (ss *SessionService) cleanupExpiredSessions() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(1)
	
	go func() {
		defer wg.Done()
		err := ss.sessionRepo.CleanExpiredSessions()
		if err != nil {
			ss.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
				Message: fmt.Sprintf("Error limpiando sesiones expiradas: %v", err),
			}))
		}
	}()

	wg.Wait()
}

// Stop detiene el servicio de sesiones
func (ss *SessionService) Stop() {
	close(ss.cleanupChannel)
	ss.cleanupWG.Wait()
}

// generateToken genera un token aleatorio seguro
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ErrSessionExpired error para sesión expirada
var ErrSessionExpired = &SessionError{Message: "sesión expirada"}

// SessionError representa un error de sesión
type SessionError struct {
	Message string
}

func (e *SessionError) Error() string {
	return e.Message
}
