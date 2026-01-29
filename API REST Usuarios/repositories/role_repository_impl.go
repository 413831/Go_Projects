package repositories

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"api-rest-usuarios/models"
)

// MockRoleRepositoryImpl implementa RoleRepository
type MockRoleRepositoryImpl struct {
	roles     map[int64]*models.Role
	userRoles map[string]*models.UserRole
	mu        sync.RWMutex
	nextID    int64
}

// NewMockRoleRepositoryImpl crea un nuevo repositorio mock de roles
func NewMockRoleRepositoryImpl() RoleRepository {
	repo := &MockRoleRepositoryImpl{
		roles:     make(map[int64]*models.Role),
		userRoles: make(map[string]*models.UserRole),
		nextID:    1,
	}

	// Crear algunos roles por defecto
	repo.roles[1] = &models.Role{ID: 1, Name: "admin", Description: "Administrador", CreatedAt: time.Now()}
	repo.roles[2] = &models.Role{ID: 2, Name: "user", Description: "Usuario", CreatedAt: time.Now()}
	repo.nextID = 3

	return repo
}

// NewPostgresRoleRepository crea un repositorio de roles para PostgreSQL.
// Implementación pendiente: por ahora retorna un mock.
func NewPostgresRoleRepository(db interface{}) RoleRepository {
	return NewMockRoleRepositoryImpl()
}

func (r *MockRoleRepositoryImpl) CreateRole(role *models.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	role.ID = r.nextID
	r.nextID++
	role.CreatedAt = time.Now()
	r.roles[role.ID] = role
	return nil
}

func (r *MockRoleRepositoryImpl) GetRoleByName(name string) (*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, role := range r.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, errors.New("rol no encontrado")
}

func (r *MockRoleRepositoryImpl) GetRoleByID(id int64) (*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, ok := r.roles[id]
	if !ok {
		return nil, errors.New("rol no encontrado")
	}
	return role, nil
}

func (r *MockRoleRepositoryImpl) GetAllRoles() ([]*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var roles []*models.Role
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *MockRoleRepositoryImpl) GrantRole(userID, roleID, grantedBy int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%d-%d", userID, roleID)
	r.userRoles[key] = &models.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		GrantedAt: time.Now(),
		GrantedBy: grantedBy,
	}
	return nil
}

func (r *MockRoleRepositoryImpl) RevokeRole(userID, roleID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%d-%d", userID, roleID)
	delete(r.userRoles, key)
	return nil
}

func (r *MockRoleRepositoryImpl) GetUserRoles(userID int64) ([]*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var roles []*models.Role
	for _, ur := range r.userRoles {
		if ur.UserID == userID {
			if role, ok := r.roles[ur.RoleID]; ok {
				roles = append(roles, role)
			}
		}
	}
	return roles, nil
}
