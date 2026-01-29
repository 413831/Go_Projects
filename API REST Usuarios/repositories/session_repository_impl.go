package repositories

import (
	"errors"
	"sync"
	"time"

	"api-rest-usuarios/models"
)

// MockSessionRepositoryImpl implementa SessionRepository
type MockSessionRepositoryImpl struct {
	sessions map[string]*models.Session
	mu       sync.RWMutex
	nextID   int64
}

// NewMockSessionRepositoryImpl crea un nuevo repositorio mock de sesiones
func NewMockSessionRepositoryImpl() SessionRepository {
	return &MockSessionRepositoryImpl{
		sessions: make(map[string]*models.Session),
		nextID:   1,
	}
}

// NewPostgresSessionRepository crea un repositorio de sesiones para PostgreSQL.
// Implementación pendiente: por ahora retorna un mock.
func NewPostgresSessionRepository(db interface{}) SessionRepository {
	return NewMockSessionRepositoryImpl()
}

func (s *MockSessionRepositoryImpl) CreateSession(session *models.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.ID = s.nextID
	s.nextID++
	session.CreatedAt = time.Now()
	session.LastActivity = time.Now()
	s.sessions[session.Token] = session
	return nil
}

func (s *MockSessionRepositoryImpl) GetSessionByToken(token string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[token]
	if !ok || !session.Active {
		return nil, errors.New("sesión no encontrada")
	}
	return session, nil
}

func (s *MockSessionRepositoryImpl) GetSessionByID(id int64) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, session := range s.sessions {
		if session.ID == id && session.Active {
			return session, nil
		}
	}
	return nil, errors.New("sesión no encontrada")
}

func (s *MockSessionRepositoryImpl) UpdateSessionActivity(sessionID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.ID == sessionID {
			session.LastActivity = time.Now()
			return nil
		}
	}
	return errors.New("sesión no encontrada")
}

func (s *MockSessionRepositoryImpl) DeactivateSession(sessionID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.ID == sessionID {
			session.Active = false
			return nil
		}
	}
	return errors.New("sesión no encontrada")
}

func (s *MockSessionRepositoryImpl) GetUserSessions(userID int64) ([]*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sessions []*models.Session
	for _, session := range s.sessions {
		if session.UserID == userID && session.Active {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (s *MockSessionRepositoryImpl) CleanExpiredSessions() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, session := range s.sessions {
		if session.ExpiresAt.Before(now) && session.Active {
			session.Active = false
			delete(s.sessions, token)
		}
	}
	return nil
}
