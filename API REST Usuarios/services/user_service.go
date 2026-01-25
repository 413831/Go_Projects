package services

import (
	"errors"
	"fmt"
	"sync"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
	"api-rest-usuarios/utils"
)

// UserService maneja la lógica de negocio de usuarios
type UserService struct {
	repo              repositories.UserRepository
	encryptionService *utils.EncryptionService
	logger            *utils.Logger
	config            interface {
		GetBCryptCost() int
	}
	mu sync.RWMutex
}

// ConfigProvider define la interfaz para obtener configuración
type ConfigProvider interface {
	GetBCryptCost() int
}

// NewUserService crea un nuevo servicio de usuarios
func NewUserService(repo repositories.UserRepository, encryptionService *utils.EncryptionService, logger *utils.Logger, config ConfigProvider) *UserService {
	return &UserService{
		repo:              repo,
		encryptionService: encryptionService,
		logger:            logger,
		config:            config,
	}
}

// Create crea un nuevo usuario
func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	s.mu.Lock()

	// Validar que el username no exista
	existingUser, _ := s.repo.GetByUsername(req.Username)
	if existingUser != nil {
		s.logger.Warn("Intento de crear usuario con username existente: " + req.Username)
		return nil, errors.New("el nombre de usuario ya existe")
	}

	// Validar que el email no exista
	existingUser, _ = s.repo.GetByEmail(req.Email)
	if existingUser != nil {
		s.logger.Warn("Intento de crear usuario con email existente: " + req.Email)
		return nil, errors.New("el email ya existe")
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(req.Password, s.config.GetBCryptCost())
	if err != nil {
		s.logger.Error("Error al hashear contraseña: " + err.Error())
		return nil, errors.New("error al procesar contraseña")
	}

	// Crear usuario
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Active:   true,
		Deleted:  false,
	}

	err = s.repo.Create(user)
	if err != nil {
		s.logger.Error("Error al crear usuario: " + err.Error())
		return nil, errors.New("error al crear usuario")
	}

	// Asignar roles si se proporcionaron
	if len(req.Roles) > 0 {
		var wg sync.WaitGroup
		errorsChan := make(chan error, len(req.Roles))

		for _, roleName := range req.Roles {
			wg.Add(1)
			go func(name string) {
				defer wg.Done()
				if err := s.grantRoleByName(user.ID, name, user.ID); err != nil {
					errorsChan <- err
				}
			}(roleName)
		}

		wg.Wait()
		close(errorsChan)

		// Verificar si hubo errores
		for err := range errorsChan {
			if err != nil {
				s.logger.Warn("Error al asignar rol: " + err.Error())
			}
		}
	}

	// Guardar PII si se proporcionó
	if req.PII != nil {
		req.PII.UserID = user.ID
		if err := s.SavePII(req.PII); err != nil {
			s.logger.Warn("Error al guardar PII: " + err.Error())
		}
	}

	// Liberar el lock antes de llamar a GetByID para evitar deadlock
	s.mu.Unlock()
	
	// Cargar datos completos (sin lock para evitar deadlock)
	result, err := s.GetByID(user.ID)
	return result, err
}

// GetByID obtiene un usuario por ID
func (s *UserService) GetByID(id int64) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cargar roles y PII en paralelo
	var wg sync.WaitGroup
	var rolesErr, piiErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		roles, err := s.repo.GetUserRoles(id)
		if err != nil {
			rolesErr = err
		} else {
			user.Roles = roles
		}
	}()

	go func() {
		defer wg.Done()
		user.PII, piiErr = s.repo.GetPIIByUserID(id)
		if user.PII != nil && piiErr == nil {
			// Desencriptar PII
			s.decryptPII(user.PII)
		}
	}()

	wg.Wait()

	if rolesErr != nil {
		s.logger.Warn("Error al cargar roles: " + rolesErr.Error())
	}
	if piiErr != nil && user.PII == nil {
		s.logger.Warn("Error al cargar PII: " + piiErr.Error())
	}

	return user, nil
}

// GetAll obtiene todos los usuarios con paginación
func (s *UserService) GetAll(page, pageSize int) ([]*models.User, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, err := s.repo.GetAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	// Cargar roles y PII para cada usuario en paralelo
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u *models.User) {
			defer wg.Done()
			roles, _ := s.repo.GetUserRoles(u.ID)
			u.Roles = roles
		}(user)
	}
	wg.Wait()

	return users, total, nil
}

// Update actualiza un usuario
func (s *UserService) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	s.mu.Lock()

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos si se proporcionaron
	if req.Username != nil {
		// Validar que el nuevo username no exista
		existingUser, _ := s.repo.GetByUsername(*req.Username)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("el nombre de usuario ya existe")
		}
		user.Username = *req.Username
	}

	if req.Email != nil {
		// Validar que el nuevo email no exista
		existingUser, _ := s.repo.GetByEmail(*req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("el email ya existe")
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password, s.config.GetBCryptCost())
		if err != nil {
			return nil, errors.New("error al procesar contraseña")
		}
		user.Password = hashedPassword
	}

	if req.Active != nil {
		user.Active = *req.Active
	}

	err = s.repo.Update(user)
	if err != nil {
		s.logger.Error("Error al actualizar usuario: " + err.Error())
		return nil, err
	}

	// Actualizar PII si se proporcionó
	if req.PII != nil {
		req.PII.UserID = id
		if err := s.SavePII(req.PII); err != nil {
			s.logger.Warn("Error al actualizar PII: " + err.Error())
		}
	}

	s.logger.Info(fmt.Sprintf("Usuario actualizado ID: %d", id))
	
	// Liberar el lock antes de llamar a GetByID para evitar deadlock
	s.mu.Unlock()
	
	// Cargar datos completos (sin lock para evitar deadlock)
	result, err := s.GetByID(id)
	return result, err
}

// Delete realiza un borrado lógico
func (s *UserService) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar usuario: " + err.Error())
		return err
	}

	s.logger.Info(fmt.Sprintf("Usuario eliminado (lógico) ID: %d", id))
	return nil
}

// GrantRole otorga un rol a un usuario
func (s *UserService) GrantRole(userID, roleID, grantedBy int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.repo.GrantRole(userID, roleID, grantedBy)
	if err != nil {
		s.logger.Error("Error al otorgar rol: " + err.Error())
		return err
	}

	s.logger.Info(fmt.Sprintf("Rol otorgado - Usuario ID: %d, Rol ID: %d", userID, roleID))
	return nil
}

// grantRoleByName otorga un rol por nombre (método auxiliar)
func (s *UserService) grantRoleByName(userID int64, roleName string, grantedBy int64) error {
	role, err := s.repo.GetRoleByName(roleName)
	if err != nil {
		// Si el rol no existe, crearlo
		role = &models.Role{
			Name:        roleName,
			Description: "Rol creado automáticamente",
		}
		err = s.repo.CreateRole(role)
		if err != nil {
			return err
		}
	}

	return s.repo.GrantRole(userID, role.ID, grantedBy)
}

// RevokeRole revoca un rol de un usuario
func (s *UserService) RevokeRole(userID, roleID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.repo.RevokeRole(userID, roleID)
	if err != nil {
		s.logger.Error("Error al revocar rol: " + err.Error())
		return err
	}

	s.logger.Info(fmt.Sprintf("Rol revocado - Usuario ID: %d, Rol ID: %d", userID, roleID))
	return nil
}

// SavePII guarda o actualiza datos PII encriptados
func (s *UserService) SavePII(pii *models.PII) error {
	// Encriptar campos sensibles
	var err error
	if pii.FirstName != "" {
		pii.FirstName, err = s.encryptionService.Encrypt(pii.FirstName)
		if err != nil {
			return err
		}
	}
	if pii.LastName != "" {
		pii.LastName, err = s.encryptionService.Encrypt(pii.LastName)
		if err != nil {
			return err
		}
	}
	if pii.PhoneNumber != "" {
		pii.PhoneNumber, err = s.encryptionService.Encrypt(pii.PhoneNumber)
		if err != nil {
			return err
		}
	}
	if pii.Address != "" {
		pii.Address, err = s.encryptionService.Encrypt(pii.Address)
		if err != nil {
			return err
		}
	}
	if pii.SSN != "" {
		pii.SSN, err = s.encryptionService.Encrypt(pii.SSN)
		if err != nil {
			return err
		}
	}

	err = s.repo.SavePII(pii)
	if err != nil {
		s.logger.Error("Error al guardar PII: " + err.Error())
		return err
	}

	s.logger.Info(fmt.Sprintf("PII guardado para usuario ID: %d", pii.UserID))
	return nil
}

// decryptPII desencripta los campos PII
func (s *UserService) decryptPII(pii *models.PII) {
	if pii.FirstName != "" {
		if decrypted, err := s.encryptionService.Decrypt(pii.FirstName); err == nil {
			pii.FirstName = decrypted
		}
	}
	if pii.LastName != "" {
		if decrypted, err := s.encryptionService.Decrypt(pii.LastName); err == nil {
			pii.LastName = decrypted
		}
	}
	if pii.PhoneNumber != "" {
		if decrypted, err := s.encryptionService.Decrypt(pii.PhoneNumber); err == nil {
			pii.PhoneNumber = decrypted
		}
	}
	if pii.Address != "" {
		if decrypted, err := s.encryptionService.Decrypt(pii.Address); err == nil {
			pii.Address = decrypted
		}
	}
	if pii.SSN != "" {
		if decrypted, err := s.encryptionService.Decrypt(pii.SSN); err == nil {
			pii.SSN = decrypted
		}
	}
}
