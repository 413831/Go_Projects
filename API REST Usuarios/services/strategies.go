package services

import (
	"errors"

	"api-rest-usuarios/models"
)

// PasswordHasher define la interfaz para hashing de contraseñas
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

// PIIEncryptor define la interfaz para encriptación de datos PII
type PIIEncryptor interface {
	Encrypt(data string) (string, error)
	Decrypt(encryptedData string) (string, error)
}

// BCryptPasswordHasher implementa PasswordHasher con bcrypt
type BCryptPasswordHasher struct {
	cost int
}

// NewBCryptPasswordHasher crea un nuevo hasher bcrypt
func NewBCryptPasswordHasher(cost int) PasswordHasher {
	return &BCryptPasswordHasher{cost: cost}
}

func (b *BCryptPasswordHasher) Hash(password string) (string, error) {
	// Importar utils.HashPassword o implementar directamente
	// Por ahora usamos la implementación existente
	return hashPasswordBcrypt(password, b.cost)
}

func (b *BCryptPasswordHasher) Verify(password, hash string) bool {
	return checkPasswordBcrypt(password, hash)
}

// AESEncryptor implementa PIIEncryptor con AES-256-GCM
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor crea un nuevo encriptador AES
func NewAESEncryptor(key string) (PIIEncryptor, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return nil, errors.New("la clave de encriptación debe tener 32 bytes")
	}
	return &AESEncryptor{key: keyBytes}, nil
}

func (a *AESEncryptor) Encrypt(data string) (string, error) {
	return encryptAES(data, a.key)
}

func (a *AESEncryptor) Decrypt(encryptedData string) (string, error) {
	return decryptAES(encryptedData, a.key)
}

// Validator define la interfaz para validaciones
type Validator interface {
	ValidateUser(user *models.CreateUserRequest) error
	ValidateUpdate(user *models.UpdateUserRequest) error
}

// UserValidator implementa validaciones de usuario
type UserValidator struct{}

// NewUserValidator crea un nuevo validador de usuarios
func NewUserValidator() Validator {
	return &UserValidator{}
}

func (v *UserValidator) ValidateUser(user *models.CreateUserRequest) error {
	if user.Username == "" {
		return errors.New("username es requerido")
	}
	if user.Email == "" {
		return errors.New("email es requerido")
	}
	if user.Password == "" {
		return errors.New("password es requerido")
	}
	if len(user.Password) < 8 {
		return errors.New("password debe tener al menos 8 caracteres")
	}
	return nil
}

func (v *UserValidator) ValidateUpdate(user *models.UpdateUserRequest) error {
	// Validaciones para actualización
	if user.Username != nil && len(*user.Username) < 3 {
		return errors.New("username debe tener al menos 3 caracteres")
	}
	if user.Email != nil && *user.Email == "" {
		return errors.New("email no puede estar vacío")
	}
	if user.Password != nil && len(*user.Password) < 8 {
		return errors.New("password debe tener al menos 8 caracteres")
	}
	return nil
}
