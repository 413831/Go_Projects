package services

import (
	"api-rest-usuarios/models"
	"api-rest-usuarios/repositories"
)

// UserQueryBuilder implementa Builder Pattern para consultas de usuarios
type UserQueryBuilder struct {
	userRepo repositories.UserRepository
	limit    int
	offset   int
	active   *bool
	deleted  *bool
	username string
	email    string
}

// NewUserQueryBuilder crea un nuevo constructor de consultas
func NewUserQueryBuilder(userRepo repositories.UserRepository) *UserQueryBuilder {
	return &UserQueryBuilder{
		userRepo: userRepo,
		limit:    10, // valor por defecto
		offset:   0,
	}
}

// Limit establece el límite de resultados
func (b *UserQueryBuilder) Limit(limit int) *UserQueryBuilder {
	b.limit = limit
	return b
}

// Offset establece el desplazamiento
func (b *UserQueryBuilder) Offset(offset int) *UserQueryBuilder {
	b.offset = offset
	return b
}

// Active filtra por usuarios activos
func (b *UserQueryBuilder) Active(active bool) *UserQueryBuilder {
	b.active = &active
	return b
}

// Deleted filtra por usuarios eliminados
func (b *UserQueryBuilder) Deleted(deleted bool) *UserQueryBuilder {
	b.deleted = &deleted
	return b
}

// Username filtra por nombre de usuario (contiene)
func (b *UserQueryBuilder) Username(username string) *UserQueryBuilder {
	b.username = username
	return b
}

// Email filtra por email (contiene)
func (b *UserQueryBuilder) Email(email string) *UserQueryBuilder {
	b.email = email
	return b
}

// Execute ejecuta la consulta y retorna los resultados
func (b *UserQueryBuilder) Execute() ([]*models.User, int, error) {
	// Para el mock, usamos GetAll y filtramos en memoria
	// En una implementación real, esto se traduciría a SQL
	
	users, err := b.userRepo.GetAll(b.limit, b.offset)
	if err != nil {
		return nil, 0, err
	}

	// Filtrar resultados
	var filtered []*models.User
	for _, user := range users {
		// Filtrar por active
		if b.active != nil && user.Active != *b.active {
			continue
		}
		
		// Filtrar por deleted
		if b.deleted != nil && user.Deleted != *b.deleted {
			continue
		}
		
		// Filtrar por username
		if b.username != "" {
			// Para simplificar, usamos coincidencia exacta
			// En una implementación real usaría LIKE
			if user.Username != b.username {
				continue
			}
		}
		
		// Filtrar por email
		if b.email != "" {
			if user.Email != b.email {
				continue
			}
		}
		
		filtered = append(filtered, user)
	}

	// Obtener total
	total, err := b.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return filtered, total, nil
}

// UserResponseBuilder implementa Builder Pattern para respuestas de usuario
type UserResponseBuilder struct {
	user     *models.User
	roles    []*models.Role
	pii      *models.PII
	sessions []*models.Session
	includePII      bool
	includeRoles    bool
	includeSessions bool
}

// NewUserResponseBuilder crea un nuevo constructor de respuesta
func NewUserResponseBuilder(user *models.User) *UserResponseBuilder {
	return &UserResponseBuilder{
		user:            user,
		includePII:      false,
		includeRoles:    false,
		includeSessions: false,
	}
}

// WithRoles incluye los roles del usuario
func (b *UserResponseBuilder) WithRoles(roles []*models.Role) *UserResponseBuilder {
	b.roles = roles
	b.includeRoles = true
	return b
}

// WithPII incluye los datos PII del usuario
func (b *UserResponseBuilder) WithPII(pii *models.PII) *UserResponseBuilder {
	b.pii = pii
	b.includePII = true
	return b
}

// WithSessions incluye las sesiones del usuario
func (b *UserResponseBuilder) WithSessions(sessions []*models.Session) *UserResponseBuilder {
	b.sessions = sessions
	b.includeSessions = true
	return b
}

// Build construye la respuesta final
func (b *UserResponseBuilder) Build() map[string]interface{} {
	response := map[string]interface{}{
		"id":         b.user.ID,
		"username":   b.user.Username,
		"email":      b.user.Email,
		"active":     b.user.Active,
		"created_at": b.user.CreatedAt,
		"updated_at": b.user.UpdatedAt,
	}

	if b.includeRoles {
		response["roles"] = b.roles
	}

	if b.includePII && b.pii != nil {
		response["pii"] = b.pii
	}

	if b.includeSessions {
		response["sessions"] = b.sessions
	}

	return response
}

// PaginatedResponseBuilder implementa Builder para respuestas paginadas
type PaginatedResponseBuilder struct {
	data       interface{}
	page       int
	pageSize   int
	total      int
	totalPages int
}

// NewPaginatedResponseBuilder crea un nuevo constructor de respuesta paginada
func NewPaginatedResponseBuilder(data interface{}, page, pageSize, total int) *PaginatedResponseBuilder {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginatedResponseBuilder{
		data:       data,
		page:       page,
		pageSize:   pageSize,
		total:      total,
		totalPages: totalPages,
	}
}

// Build construye la respuesta paginada
func (b *PaginatedResponseBuilder) Build() map[string]interface{} {
	return map[string]interface{}{
		"data":        b.data,
		"page":        b.page,
		"page_size":   b.pageSize,
		"total":       b.total,
		"total_pages": b.totalPages,
	}
}
