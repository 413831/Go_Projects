package repositories

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"api-rest-usuarios/models"
)

// MockUserRepository es una implementación en memoria para pruebas
type MockUserRepository struct {
	users    map[int64]*models.User
	roles    map[int64]*models.Role
	userRoles map[string]*models.UserRole
	pii      map[int64]*models.PII
	sessions map[string]*models.Session
	mu       sync.RWMutex
	nextID   int64
}

// NewMockUserRepository crea un nuevo repositorio mock
func NewMockUserRepository() UserRepository {
	repo := &MockUserRepository{
		users:     make(map[int64]*models.User),
		roles:     make(map[int64]*models.Role),
		userRoles: make(map[string]*models.UserRole),
		pii:       make(map[int64]*models.PII),
		sessions:  make(map[string]*models.Session),
		nextID:    1,
	}

	// Crear algunos roles por defecto
	repo.roles[1] = &models.Role{ID: 1, Name: "admin", Description: "Administrador", CreatedAt: time.Now()}
	repo.roles[2] = &models.Role{ID: 2, Name: "user", Description: "Usuario", CreatedAt: time.Now()}
	repo.nextID = 3

	return repo
}

func (m *MockUserRepository) Create(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.ID = m.nextID
	m.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(id int64) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]
	if !ok || user.Deleted {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == username && !user.Deleted {
			return user, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email && !user.Deleted {
			return user, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *MockUserRepository) GetAll(limit, offset int) ([]*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var users []*models.User
	count := 0
	for _, user := range m.users {
		if !user.Deleted {
			if count >= offset {
				users = append(users, user)
				if len(users) >= limit {
					break
				}
			}
			count++
		}
	}
	return users, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[user.ID]; !ok {
		return errors.New("usuario no encontrado")
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[id]
	if !ok {
		return errors.New("usuario no encontrado")
	}
	user.Deleted = true
	user.UpdatedAt = time.Now()
	return nil
}

func (m *MockUserRepository) Count() (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, user := range m.users {
		if !user.Deleted {
			count++
		}
	}
	return count, nil
}

func (m *MockUserRepository) GrantRole(userID, roleID, grantedBy int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%d-%d", userID, roleID)
	m.userRoles[key] = &models.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		GrantedAt: time.Now(),
		GrantedBy: grantedBy,
	}
	return nil
}

func (m *MockUserRepository) RevokeRole(userID, roleID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%d-%d", userID, roleID)
	delete(m.userRoles, key)
	return nil
}

func (m *MockUserRepository) GetUserRoles(userID int64) ([]*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var roles []*models.Role
	for _, ur := range m.userRoles {
		if ur.UserID == userID {
			if role, ok := m.roles[ur.RoleID]; ok {
				roles = append(roles, role)
			}
		}
	}
	return roles, nil
}

func (m *MockUserRepository) GetRoleByName(name string) (*models.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, errors.New("rol no encontrado")
}

func (m *MockUserRepository) CreateRole(role *models.Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	role.ID = m.nextID
	m.nextID++
	role.CreatedAt = time.Now()
	m.roles[role.ID] = role
	return nil
}

func (m *MockUserRepository) SavePII(pii *models.PII) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if pii.ID == 0 {
		pii.ID = m.nextID
		m.nextID++
		pii.CreatedAt = time.Now()
	}
	pii.UpdatedAt = time.Now()
	m.pii[pii.UserID] = pii
	return nil
}

func (m *MockUserRepository) GetPIIByUserID(userID int64) (*models.PII, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pii, ok := m.pii[userID]
	if !ok {
		return nil, nil
	}
	return pii, nil
}

func (m *MockUserRepository) UpdatePII(pii *models.PII) error {
	return m.SavePII(pii)
}

func (m *MockUserRepository) CreateSession(session *models.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session.ID = m.nextID
	m.nextID++
	session.CreatedAt = time.Now()
	session.LastActivity = time.Now()
	m.sessions[session.Token] = session
	return nil
}

func (m *MockUserRepository) GetSessionByToken(token string) (*models.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[token]
	if !ok || !session.Active {
		return nil, errors.New("sesión no encontrada")
	}
	return session, nil
}

func (m *MockUserRepository) UpdateSessionActivity(sessionID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, session := range m.sessions {
		if session.ID == sessionID {
			session.LastActivity = time.Now()
			return nil
		}
	}
	return errors.New("sesión no encontrada")
}

func (m *MockUserRepository) DeactivateSession(sessionID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, session := range m.sessions {
		if session.ID == sessionID {
			session.Active = false
			return nil
		}
	}
	return errors.New("sesión no encontrada")
}

func (m *MockUserRepository) GetUserSessions(userID int64) ([]*models.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sessions []*models.Session
	for _, session := range m.sessions {
		if session.UserID == userID && session.Active {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (m *MockUserRepository) CleanExpiredSessions() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if session.ExpiresAt.Before(now) && session.Active {
			session.Active = false
			delete(m.sessions, token)
		}
	}
	return nil
}
