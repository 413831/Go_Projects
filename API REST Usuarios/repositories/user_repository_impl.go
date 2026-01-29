package repositories

import (
	"errors"
	"sync"
	"time"

	"api-rest-usuarios/models"
)

// MockUserRepositoryImpl implementa UserRepository
type MockUserRepositoryImpl struct {
	users  map[int64]*models.User
	mu     sync.RWMutex
	nextID int64
}

// NewMockUserRepositoryImpl crea un nuevo repositorio mock de usuarios
func NewMockUserRepositoryImpl() UserRepository {
	return &MockUserRepositoryImpl{
		users:  make(map[int64]*models.User),
		nextID: 1,
	}
}

// NewPostgresUserRepository crea un repositorio de usuarios para PostgreSQL.
// Implementación pendiente: por ahora retorna un mock.
func NewPostgresUserRepository(db interface{}) UserRepository {
	return NewMockUserRepositoryImpl()
}

func (m *MockUserRepositoryImpl) Create(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.ID = m.nextID
	m.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepositoryImpl) GetByID(id int64) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]
	if !ok || user.Deleted {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

func (m *MockUserRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == username && !user.Deleted {
			return user, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *MockUserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email && !user.Deleted {
			return user, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *MockUserRepositoryImpl) GetAll(limit, offset int) ([]*models.User, error) {
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

func (m *MockUserRepositoryImpl) Update(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[user.ID]; !ok {
		return errors.New("usuario no encontrado")
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepositoryImpl) Delete(id int64) error {
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

func (m *MockUserRepositoryImpl) Count() (int, error) {
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
