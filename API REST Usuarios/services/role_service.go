package services

import (
	"fmt"
	"sync"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

// RoleServiceInterface define la interfaz para el servicio de roles
type RoleServiceInterface interface {
	CreateRole(role *CreateRoleRequest) (*models.Role, error)
	GetRoleByID(id int64) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() ([]*models.Role, error)
	UpdateRole(id int64, req *UpdateRoleRequest) (*models.Role, error)
	DeleteRole(id int64) error
}

// CreateRoleRequest representa la solicitud para crear un rol
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"required,max=200"`
}

// UpdateRoleRequest representa la solicitud para actualizar un rol
type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=200"`
}

// RoleService implementa RoleServiceInterface
type RoleService struct {
	roleRepo       repositories.RoleRepository
	eventPublisher EventPublisher
	mu             sync.RWMutex
}

// NewRoleService crea una nueva instancia del servicio de roles
func NewRoleService(
	roleRepo repositories.RoleRepository,
	eventPublisher EventPublisher,
) RoleServiceInterface {
	return &RoleService{
		roleRepo:       roleRepo,
		eventPublisher: eventPublisher,
	}
}

func (rs *RoleService) CreateRole(req *CreateRoleRequest) (*models.Role, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Verificar que el rol no exista
	existingRole, _ := rs.roleRepo.GetRoleByName(req.Name)
	if existingRole != nil {
		return nil, fmt.Errorf("el rol %s ya existe", req.Name)
	}

	// Crear rol
	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := rs.roleRepo.CreateRole(role); err != nil {
		rs.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			Message: fmt.Sprintf("Error creando rol: %v", err),
		}))
		return nil, fmt.Errorf("error creando rol: %w", err)
	}

	// Publicar evento de rol creado
	rs.eventPublisher.Publish(NewUserEvent("role.created", map[string]interface{}{
		"role_id":   role.ID,
		"role_name": role.Name,
	}))

	return role, nil
}

func (rs *RoleService) GetRoleByID(id int64) (*models.Role, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	return rs.roleRepo.GetRoleByID(id)
}

func (rs *RoleService) GetRoleByName(name string) (*models.Role, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	return rs.roleRepo.GetRoleByName(name)
}

func (rs *RoleService) GetAllRoles() ([]*models.Role, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	return rs.roleRepo.GetAllRoles()
}

func (rs *RoleService) UpdateRole(id int64, req *UpdateRoleRequest) (*models.Role, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Obtener rol existente
	role, err := rs.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("rol no encontrado: %w", err)
	}

	// Actualizar campos si se proporcionaron
	if req.Name != nil {
		// Verificar que el nuevo nombre no exista
		existingRole, _ := rs.roleRepo.GetRoleByName(*req.Name)
		if existingRole != nil && existingRole.ID != id {
			return nil, fmt.Errorf("el rol %s ya existe", *req.Name)
		}
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	// Nota: En una implementación real necesitaríamos un método UpdateRole en el repositorio
	// Por ahora retornamos el rol actualizado
	
	// Publicar evento de rol actualizado
	rs.eventPublisher.Publish(NewUserEvent("role.updated", map[string]interface{}{
		"role_id":   role.ID,
		"role_name": role.Name,
	}))

	return role, nil
}

func (rs *RoleService) DeleteRole(id int64) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Obtener rol para evento
	role, err := rs.roleRepo.GetRoleByID(id)
	if err != nil {
		return fmt.Errorf("rol no encontrado: %w", err)
	}

	// Nota: En una implementación real necesitaríamos un método DeleteRole en el repositorio
	// Por ahora solo publicamos el evento
	
	// Publicar evento de rol eliminado
	rs.eventPublisher.Publish(NewUserEvent("role.deleted", map[string]interface{}{
		"role_id":   role.ID,
		"role_name": role.Name,
	}))

	return nil
}
