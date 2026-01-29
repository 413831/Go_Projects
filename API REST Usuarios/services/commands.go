package services

import (
	"errors"

	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

// Command define la interfaz base para comandos
type Command interface {
	Execute() error
	Undo() error
}

// CommandHandler define la interfaz para manejar comandos
type CommandHandler interface {
	Handle(cmd Command) error
}

// CreateUserCommand implementa Command para crear usuarios
type CreateUserCommand struct {
	userRepo   repositories.UserRepository
	roleRepo   repositories.RoleRepository
	piiRepo    repositories.PIIRepository
	hasher     PasswordHasher
	encryptor  PIIEncryptor
	validator  Validator
	request    *models.CreateUserRequest
	result     *models.User
}

// NewCreateUserCommand crea un nuevo comando para crear usuario
func NewCreateUserCommand(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	piiRepo repositories.PIIRepository,
	hasher PasswordHasher,
	encryptor PIIEncryptor,
	validator Validator,
	request *models.CreateUserRequest,
) Command {
	return &CreateUserCommand{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		piiRepo:   piiRepo,
		hasher:    hasher,
		encryptor: encryptor,
		validator: validator,
		request:   request,
	}
}

func (c *CreateUserCommand) Execute() error {
	// Validar solicitud
	if err := c.validator.ValidateUser(c.request); err != nil {
		return err
	}

	// Verificar que el username no exista
	existingUser, _ := c.userRepo.GetByUsername(c.request.Username)
	if existingUser != nil {
		return errors.New("el nombre de usuario ya existe")
	}

	// Verificar que el email no exista
	existingUser, _ = c.userRepo.GetByEmail(c.request.Email)
	if existingUser != nil {
		return errors.New("el email ya existe")
	}

	// Hash de la contraseña
	hashedPassword, err := c.hasher.Hash(c.request.Password)
	if err != nil {
		return errors.New("error al procesar contraseña")
	}

	// Crear usuario
	user := &models.User{
		Username: c.request.Username,
		Email:    c.request.Email,
		Password: hashedPassword,
		Active:   true,
		Deleted:  false,
	}

	if err := c.userRepo.Create(user); err != nil {
		return errors.New("error al crear usuario")
	}

	// Asignar roles si se proporcionaron
	if len(c.request.Roles) > 0 {
		for _, roleName := range c.request.Roles {
			role, err := c.roleRepo.GetRoleByName(roleName)
			if err == nil {
				c.roleRepo.GrantRole(user.ID, role.ID, user.ID)
			}
		}
	}

	// Guardar PII si se proporcionó
	if c.request.PII != nil {
		c.request.PII.UserID = user.ID
		
		// Encriptar datos PII
		if encryptedFirstName, err := c.encryptor.Encrypt(c.request.PII.FirstName); err == nil {
			c.request.PII.FirstName = encryptedFirstName
		}
		if encryptedLastName, err := c.encryptor.Encrypt(c.request.PII.LastName); err == nil {
			c.request.PII.LastName = encryptedLastName
		}
		if encryptedPhone, err := c.encryptor.Encrypt(c.request.PII.PhoneNumber); err == nil {
			c.request.PII.PhoneNumber = encryptedPhone
		}
		if encryptedSSN, err := c.encryptor.Encrypt(c.request.PII.SSN); err == nil {
			c.request.PII.SSN = encryptedSSN
		}
		
		c.piiRepo.SavePII(c.request.PII)
	}

	c.result = user
	return nil
}

func (c *CreateUserCommand) Undo() error {
	if c.result != nil {
		return c.userRepo.Delete(c.result.ID)
	}
	return nil
}

// GetResult retorna el resultado del comando
func (c *CreateUserCommand) GetResult() *models.User {
	return c.result
}

// UpdateUserCommand implementa Command para actualizar usuarios
type UpdateUserCommand struct {
	userRepo   repositories.UserRepository
	piiRepo    repositories.PIIRepository
	hasher     PasswordHasher
	encryptor  PIIEncryptor
	validator  Validator
	userID     int64
	request    *models.UpdateUserRequest
	original   *models.User
}

// NewUpdateUserCommand crea un nuevo comando para actualizar usuario
func NewUpdateUserCommand(
	userRepo repositories.UserRepository,
	piiRepo repositories.PIIRepository,
	hasher PasswordHasher,
	encryptor PIIEncryptor,
	validator Validator,
	userID int64,
	request *models.UpdateUserRequest,
) Command {
	return &UpdateUserCommand{
		userRepo:  userRepo,
		piiRepo:   piiRepo,
		hasher:    hasher,
		encryptor: encryptor,
		validator: validator,
		userID:    userID,
		request:   request,
	}
}

func (c *UpdateUserCommand) Execute() error {
	// Validar solicitud
	if err := c.validator.ValidateUpdate(c.request); err != nil {
		return err
	}

	// Obtener usuario original
	user, err := c.userRepo.GetByID(c.userID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	c.original = &models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Active:    user.Active,
		Deleted:   user.Deleted,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Actualizar campos si se proporcionaron
	if c.request.Username != nil {
		user.Username = *c.request.Username
	}
	if c.request.Email != nil {
		user.Email = *c.request.Email
	}
	if c.request.Password != nil {
		hashedPassword, err := c.hasher.Hash(*c.request.Password)
		if err != nil {
			return errors.New("error al procesar contraseña")
		}
		user.Password = hashedPassword
	}
	if c.request.Active != nil {
		user.Active = *c.request.Active
	}

	if err := c.userRepo.Update(user); err != nil {
		return errors.New("error al actualizar usuario")
	}

	// Actualizar PII si se proporcionó
	if c.request.PII != nil {
		c.request.PII.UserID = user.ID
		
		// Encriptar datos PII
		if c.request.PII.FirstName != "" {
			if encrypted, err := c.encryptor.Encrypt(c.request.PII.FirstName); err == nil {
				c.request.PII.FirstName = encrypted
			}
		}
		if c.request.PII.LastName != "" {
			if encrypted, err := c.encryptor.Encrypt(c.request.PII.LastName); err == nil {
				c.request.PII.LastName = encrypted
			}
		}
		
		c.piiRepo.UpdatePII(c.request.PII)
	}

	return nil
}

func (c *UpdateUserCommand) Undo() error {
	if c.original != nil {
		return c.userRepo.Update(c.original)
	}
	return nil
}

// DeleteUserCommand implementa Command para eliminar usuarios
type DeleteUserCommand struct {
	userRepo repositories.UserRepository
	userID   int64
	deleted  bool
}

// NewDeleteUserCommand crea un nuevo comando para eliminar usuario
func NewDeleteUserCommand(userRepo repositories.UserRepository, userID int64) Command {
	return &DeleteUserCommand{
		userRepo: userRepo,
		userID:   userID,
	}
}

func (c *DeleteUserCommand) Execute() error {
	if err := c.userRepo.Delete(c.userID); err != nil {
		return err
	}
	c.deleted = true
	return nil
}

func (c *DeleteUserCommand) Undo() error {
	if c.deleted {
		// Para undo, marcar como no eliminado (necesitaría método específico)
		return errors.New("undo no implementado para delete")
	}
	return nil
}
