package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
	"api-rest-usuarios/utils"
)

// SessionService maneja la lógica de negocio de sesiones
type SessionService struct {
	repo           repositories.UserRepository
	logger         *utils.Logger
	cleanupChannel chan bool
	cleanupWG      sync.WaitGroup
	mu             sync.Mutex
}

// NewSessionService crea un nuevo servicio de sesiones
func NewSessionService(repo repositories.UserRepository, logger *utils.Logger) *SessionService {
	ss := &SessionService{
		repo:           repo,
		logger:         logger,
		cleanupChannel: make(chan bool),
	}
	// Iniciar goroutine para limpieza periódica de sesiones
	go ss.startCleanupRoutine()
	return ss
}

// CreateSession crea una nueva sesión para un usuario
func (ss *SessionService) CreateSession(userID int64, ipAddress, userAgent string) (*models.Session, error) {
	token, err := generateToken()
	if err != nil {
		ss.logger.Error("Error al generar token: " + err.Error())
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

	err = ss.repo.CreateSession(session)
	if err != nil {
		ss.logger.Error("Error al crear sesión: " + err.Error())
		return nil, err
	}

	ss.logger.Info(fmt.Sprintf("Sesión creada para usuario ID: %d", userID))
	return session, nil
}

// ValidateSession valida si una sesión es válida
func (ss *SessionService) ValidateSession(token string) (*models.Session, error) {
	session, err := ss.repo.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		ss.repo.DeactivateSession(session.ID)
		return nil, ErrSessionExpired
	}

	// Actualizar última actividad en goroutine para no bloquear
	go func() {
		ss.repo.UpdateSessionActivity(session.ID)
	}()

	return session, nil
}

// GetUserSessions obtiene todas las sesiones activas de un usuario
func (ss *SessionService) GetUserSessions(userID int64) ([]*models.Session, error) {
	return ss.repo.GetUserSessions(userID)
}

// DeactivateSession desactiva una sesión
func (ss *SessionService) DeactivateSession(sessionID int64) error {
	err := ss.repo.DeactivateSession(sessionID)
	if err != nil {
		ss.logger.Error("Error al desactivar sesión: " + err.Error())
		return err
	}
	ss.logger.Info(fmt.Sprintf("Sesión desactivada ID: %d", sessionID))
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
		err := ss.repo.CleanExpiredSessions()
		if err != nil {
			ss.logger.Error("Error al limpiar sesiones expiradas: " + err.Error())
		} else {
			ss.logger.Info("Sesiones expiradas limpiadas")
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

var ErrSessionExpired = &SessionError{Message: "sesión expirada"}

type SessionError struct {
	Message string
}

func (e *SessionError) Error() string {
	return e.Message
}
