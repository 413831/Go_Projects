package repositories

import (
	"api-rest-usuarios/models"
)

// UserRepository define la interfaz para el repositorio de usuarios
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll(limit, offset int) ([]*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error // Borrado lógico
	Count() (int, error)
}

// RoleRepository define la interfaz para el repositorio de roles
type RoleRepository interface {
	CreateRole(role *models.Role) error
	GetRoleByName(name string) (*models.Role, error)
	GetRoleByID(id int64) (*models.Role, error)
	GetAllRoles() ([]*models.Role, error)
	GrantRole(userID, roleID, grantedBy int64) error
	RevokeRole(userID, roleID int64) error
	GetUserRoles(userID int64) ([]*models.Role, error)
}

// PIIRepository define la interfaz para el repositorio de datos PII
type PIIRepository interface {
	SavePII(pii *models.PII) error
	GetPIIByUserID(userID int64) (*models.PII, error)
	UpdatePII(pii *models.PII) error
	DeletePII(userID int64) error
}

// SessionRepository define la interfaz para el repositorio de sesiones
type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	GetSessionByID(id int64) (*models.Session, error)
	UpdateSessionActivity(sessionID int64) error
	DeactivateSession(sessionID int64) error
	GetUserSessions(userID int64) ([]*models.Session, error)
	CleanExpiredSessions() error
}

// CompositeRepository combina todos los repositorios para compatibilidad
type CompositeRepository interface {
	UserRepository
	RoleRepository
	PIIRepository
	SessionRepository
}
