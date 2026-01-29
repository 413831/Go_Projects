package services

import (
	"fmt"
	"sync"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

// UserServiceInterface define la interfaz para el servicio de usuarios
type UserServiceInterface interface {
	Create(req *models.CreateUserRequest) (*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetAll(page, pageSize int) ([]*models.User, int, error)
	Update(id int64, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id int64) error
	GrantRole(userID, roleID, grantedBy int64) error
	RevokeRole(userID, roleID int64) error
	SavePII(pii *models.PII) error
	GetUserRoles(userID int64) ([]*models.Role, error)
}

// UserService implementa UserServiceInterface con patrones de diseño
type UserService struct {
	userRepo       repositories.UserRepository
	roleRepo       repositories.RoleRepository
	piiRepo        repositories.PIIRepository
	hasher         PasswordHasher
	encryptor      PIIEncryptor
	validator      Validator
	eventPublisher EventPublisher
	commandHandler CommandHandler
	mu             sync.RWMutex
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	piiRepo repositories.PIIRepository,
	hasher PasswordHasher,
	encryptor PIIEncryptor,
	validator Validator,
	eventPublisher EventPublisher,
) UserServiceInterface {
	service := &UserService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		piiRepo:        piiRepo,
		hasher:         hasher,
		encryptor:      encryptor,
		validator:      validator,
		eventPublisher: eventPublisher,
		commandHandler: NewDefaultCommandHandler(),
	}

	return service
}

// NewUserServiceV2 mantiene compatibilidad con application_v2/factory.
func NewUserServiceV2(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	piiRepo repositories.PIIRepository,
	hasher PasswordHasher,
	encryptor PIIEncryptor,
	validator Validator,
	eventPublisher EventPublisher,
) UserServiceInterface {
	return NewUserService(userRepo, roleRepo, piiRepo, hasher, encryptor, validator, eventPublisher)
}

func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	s.mu.Lock()

	// Crear y ejecutar comando
	cmd := NewCreateUserCommand(
		s.userRepo,
		s.roleRepo,
		s.piiRepo,
		s.hasher,
		s.encryptor,
		s.validator,
		req,
	)

	if err := s.commandHandler.Handle(cmd); err != nil {
		// Publicar evento de error
		s.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			Message: fmt.Sprintf("Error creando usuario: %v", err),
		}))
		return nil, err
	}

	// Obtener resultado del comando
	createCmd := cmd.(*CreateUserCommand)
	user := createCmd.GetResult()

	// Publicar evento de usuario creado
	s.eventPublisher.Publish(NewUserEvent("user.created", &UserData{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}))

	createdID := user.ID
	s.mu.Unlock()

	// Cargar datos completos
	return s.GetByID(createdID)
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cargar roles y PII en paralelo usando Builder Pattern
	var wg sync.WaitGroup
	var roles []*models.Role
	var pii *models.PII

	wg.Add(2)

	go func() {
		defer wg.Done()
		roles, _ = s.roleRepo.GetUserRoles(id)
	}()

	go func() {
		defer wg.Done()
		pii, _ = s.piiRepo.GetPIIByUserID(id)
		if pii != nil {
			s.decryptPII(pii)
		}
	}()

	wg.Wait()

	user.Roles = roles
	user.PII = pii

	return user, nil
}

func (s *UserService) GetAll(page, pageSize int) ([]*models.User, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Usar Builder Pattern para consulta
	queryBuilder := NewUserQueryBuilder(s.userRepo).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Active(true).
		Deleted(false)

	users, total, err := queryBuilder.Execute()
	if err != nil {
		return nil, 0, err
	}

	// Cargar roles para cada usuario en paralelo
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u *models.User) {
			defer wg.Done()
			roles, _ := s.roleRepo.GetUserRoles(u.ID)
			u.Roles = roles
		}(user)
	}
	wg.Wait()

	return users, total, nil
}

func (s *UserService) Update(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	s.mu.Lock()

	// Crear y ejecutar comando
	cmd := NewUpdateUserCommand(
		s.userRepo,
		s.piiRepo,
		s.hasher,
		s.encryptor,
		s.validator,
		id,
		req,
	)

	if err := s.commandHandler.Handle(cmd); err != nil {
		s.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			UserID:  id,
			Message: fmt.Sprintf("Error actualizando usuario: %v", err),
		}))
		return nil, err
	}

	// Publicar evento de usuario actualizado
	s.eventPublisher.Publish(NewUserEvent("user.updated", &UserData{
		UserID: id,
	}))

	updatedID := id
	s.mu.Unlock()

	// Cargar datos completos
	return s.GetByID(updatedID)
}

func (s *UserService) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Obtener usuario para evento
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Crear y ejecutar comando
	cmd := NewDeleteUserCommand(s.userRepo, id)
	if err := s.commandHandler.Handle(cmd); err != nil {
		s.eventPublisher.Publish(NewUserEvent("error.occurred", &ErrorData{
			UserID:  id,
			Message: fmt.Sprintf("Error eliminando usuario: %v", err),
		}))
		return err
	}

	// Publicar evento de usuario eliminado
	s.eventPublisher.Publish(NewUserEvent("user.deleted", &UserData{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}))

	return nil
}

func (s *UserService) GrantRole(userID, roleID, grantedBy int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.roleRepo.GrantRole(userID, roleID, grantedBy)
}

func (s *UserService) RevokeRole(userID, roleID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.roleRepo.RevokeRole(userID, roleID)
}

func (s *UserService) SavePII(pii *models.PII) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encriptar datos PII antes de guardar
	if pii.FirstName != "" {
		if encrypted, err := s.encryptor.Encrypt(pii.FirstName); err == nil {
			pii.FirstName = encrypted
		}
	}
	if pii.LastName != "" {
		if encrypted, err := s.encryptor.Encrypt(pii.LastName); err == nil {
			pii.LastName = encrypted
		}
	}
	if pii.PhoneNumber != "" {
		if encrypted, err := s.encryptor.Encrypt(pii.PhoneNumber); err == nil {
			pii.PhoneNumber = encrypted
		}
	}
	if pii.SSN != "" {
		if encrypted, err := s.encryptor.Encrypt(pii.SSN); err == nil {
			pii.SSN = encrypted
		}
	}

	return s.piiRepo.SavePII(pii)
}

func (s *UserService) GetUserRoles(userID int64) ([]*models.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.roleRepo.GetUserRoles(userID)
}

// decryptPII desencripta los datos PII de un usuario
func (s *UserService) decryptPII(pii *models.PII) {
	if pii.FirstName != "" {
		if decrypted, err := s.encryptor.Decrypt(pii.FirstName); err == nil {
			pii.FirstName = decrypted
		}
	}
	if pii.LastName != "" {
		if decrypted, err := s.encryptor.Decrypt(pii.LastName); err == nil {
			pii.LastName = decrypted
		}
	}
	if pii.PhoneNumber != "" {
		if decrypted, err := s.encryptor.Decrypt(pii.PhoneNumber); err == nil {
			pii.PhoneNumber = decrypted
		}
	}
	if pii.SSN != "" {
		if decrypted, err := s.encryptor.Decrypt(pii.SSN); err == nil {
			pii.SSN = decrypted
		}
	}
}

// DefaultCommandHandler implementa CommandHandler básico
type DefaultCommandHandler struct{}

// NewDefaultCommandHandler crea un nuevo manejador de comandos por defecto
func NewDefaultCommandHandler() CommandHandler {
	return &DefaultCommandHandler{}
}

func (h *DefaultCommandHandler) Handle(cmd Command) error {
	return cmd.Execute()
}
