package models

import (
	"time"
)

// User representa un usuario en el sistema
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // No se expone en JSON
	Active    bool      `json:"active" db:"active"`
	Deleted   bool      `json:"deleted" db:"deleted"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Roles     []*Role   `json:"roles,omitempty" db:"-"`
	PII       *PII      `json:"pii,omitempty" db:"-"`
}

// Role representa un rol del sistema
type Role struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserRole representa la relación entre usuario y rol
type UserRole struct {
	UserID    int64     `json:"user_id" db:"user_id"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	GrantedAt time.Time `json:"granted_at" db:"granted_at"`
	GrantedBy int64     `json:"granted_by" db:"granted_by"`
}

// PII (Personally Identifiable Information) - Datos personales sensibles
type PII struct {
	ID            int64     `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	FirstName     string    `json:"first_name" db:"first_name"`     // Encriptado
	LastName      string    `json:"last_name" db:"last_name"`       // Encriptado
	PhoneNumber   string    `json:"phone_number" db:"phone_number"` // Encriptado
	Address       string    `json:"address" db:"address"`           // Encriptado
	SSN           string    `json:"ssn" db:"ssn"`                   // Encriptado (Social Security Number)
	DateOfBirth   time.Time `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Session representa una sesión de usuario
type Session struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Token        string    `json:"token" db:"token"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Active       bool      `json:"active" db:"active"`
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Roles    []string `json:"roles,omitempty"`
	PII      *PII     `json:"pii,omitempty"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
	Active   *bool   `json:"active,omitempty"`
	PII      *PII    `json:"pii,omitempty"`
}

// GrantRoleRequest representa la solicitud para otorgar un rol
type GrantRoleRequest struct {
	RoleID int64 `json:"role_id" validate:"required"`
}
