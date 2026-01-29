package repositories

import (
	"sync"
	"time"

	"api-rest-usuarios/models"
)

// MockPIIRepositoryImpl implementa PIIRepository
type MockPIIRepositoryImpl struct {
	pii    map[int64]*models.PII
	mu     sync.RWMutex
	nextID int64
}

// NewMockPIIRepositoryImpl crea un nuevo repositorio mock de PII
func NewMockPIIRepositoryImpl() PIIRepository {
	return &MockPIIRepositoryImpl{
		pii:    make(map[int64]*models.PII),
		nextID: 1,
	}
}

// NewPostgresPIIRepository crea un repositorio de PII para PostgreSQL.
// Implementación pendiente: por ahora retorna un mock.
func NewPostgresPIIRepository(db interface{}) PIIRepository {
	return NewMockPIIRepositoryImpl()
}

func (p *MockPIIRepositoryImpl) SavePII(pii *models.PII) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pii.ID == 0 {
		pii.ID = p.nextID
		p.nextID++
		pii.CreatedAt = time.Now()
	}
	pii.UpdatedAt = time.Now()
	p.pii[pii.UserID] = pii
	return nil
}

func (p *MockPIIRepositoryImpl) GetPIIByUserID(userID int64) (*models.PII, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	pii, ok := p.pii[userID]
	if !ok {
		return nil, nil
	}
	return pii, nil
}

func (p *MockPIIRepositoryImpl) UpdatePII(pii *models.PII) error {
	return p.SavePII(pii)
}

func (p *MockPIIRepositoryImpl) DeletePII(userID int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.pii, userID)
	return nil
}
