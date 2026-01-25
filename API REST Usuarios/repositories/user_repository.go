package repositories

import (
	"database/sql"
	"errors"
	"time"

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
	
	// Roles
	GrantRole(userID, roleID, grantedBy int64) error
	RevokeRole(userID, roleID int64) error
	GetUserRoles(userID int64) ([]*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	CreateRole(role *models.Role) error
	
	// PII
	SavePII(pii *models.PII) error
	GetPIIByUserID(userID int64) (*models.PII, error)
	UpdatePII(pii *models.PII) error
	
	// Sessions
	CreateSession(session *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	UpdateSessionActivity(sessionID int64) error
	DeactivateSession(sessionID int64) error
	GetUserSessions(userID int64) ([]*models.Session, error)
	CleanExpiredSessions() error
}

// userRepository implementa UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create crea un nuevo usuario
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, active, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Active,
		user.Deleted,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)
	return err
}

// GetByID obtiene un usuario por ID
func (r *userRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password, active, deleted, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted = false
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}

// GetByUsername obtiene un usuario por nombre de usuario
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, active, deleted, created_at, updated_at
		FROM users
		WHERE username = $1 AND deleted = false
	`
	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}

// GetByEmail obtiene un usuario por email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, active, deleted, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted = false
	`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Deleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}

// GetAll obtiene todos los usuarios con paginación
func (r *userRepository) GetAll(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password, active, deleted, created_at, updated_at
		FROM users
		WHERE deleted = false
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Active,
			&user.Deleted,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Update actualiza un usuario
func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, active = $4, updated_at = $5
		WHERE id = $6 AND deleted = false
	`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Active, time.Now(), user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}

// Delete realiza un borrado lógico
func (r *userRepository) Delete(id int64) error {
	query := `UPDATE users SET deleted = true, updated_at = $1 WHERE id = $2`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}

// Count cuenta el total de usuarios activos
func (r *userRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE deleted = false`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// GrantRole otorga un rol a un usuario
func (r *userRepository) GrantRole(userID, roleID, grantedBy int64) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, granted_at, granted_by)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.Exec(query, userID, roleID, time.Now(), grantedBy)
	return err
}

// RevokeRole revoca un rol de un usuario
func (r *userRepository) RevokeRole(userID, roleID int64) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(query, userID, roleID)
	return err
}

// GetUserRoles obtiene los roles de un usuario
func (r *userRepository) GetUserRoles(userID int64) ([]*models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// GetRoleByName obtiene un rol por nombre
func (r *userRepository) GetRoleByName(name string) (*models.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE name = $1`
	role := &models.Role{}
	err := r.db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("rol no encontrado")
		}
		return nil, err
	}
	return role, nil
}

// CreateRole crea un nuevo rol
func (r *userRepository) CreateRole(role *models.Role) error {
	query := `
		INSERT INTO roles (name, description, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRow(query, role.Name, role.Description, time.Now()).Scan(&role.ID)
	return err
}

// SavePII guarda datos PII de un usuario
func (r *userRepository) SavePII(pii *models.PII) error {
	query := `
		INSERT INTO pii (user_id, first_name, last_name, phone_number, address, ssn, date_of_birth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE
		SET first_name = $2, last_name = $3, phone_number = $4, address = $5, ssn = $6, date_of_birth = $7, updated_at = $9
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		pii.UserID,
		pii.FirstName,
		pii.LastName,
		pii.PhoneNumber,
		pii.Address,
		pii.SSN,
		pii.DateOfBirth,
		time.Now(),
		time.Now(),
	).Scan(&pii.ID)
	return err
}

// GetPIIByUserID obtiene los datos PII de un usuario
func (r *userRepository) GetPIIByUserID(userID int64) (*models.PII, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone_number, address, ssn, date_of_birth, created_at, updated_at
		FROM pii
		WHERE user_id = $1
	`
	pii := &models.PII{}
	err := r.db.QueryRow(query, userID).Scan(
		&pii.ID,
		&pii.UserID,
		&pii.FirstName,
		&pii.LastName,
		&pii.PhoneNumber,
		&pii.Address,
		&pii.SSN,
		&pii.DateOfBirth,
		&pii.CreatedAt,
		&pii.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No es un error, simplemente no hay PII
		}
		return nil, err
	}
	return pii, nil
}

// UpdatePII actualiza los datos PII
func (r *userRepository) UpdatePII(pii *models.PII) error {
	query := `
		UPDATE pii
		SET first_name = $1, last_name = $2, phone_number = $3, address = $4, ssn = $5, date_of_birth = $6, updated_at = $7
		WHERE user_id = $8
	`
	_, err := r.db.Exec(query, pii.FirstName, pii.LastName, pii.PhoneNumber, pii.Address, pii.SSN, pii.DateOfBirth, time.Now(), pii.UserID)
	return err
}

// CreateSession crea una nueva sesión
func (r *userRepository) CreateSession(session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, ip_address, user_agent, created_at, expires_at, last_activity, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		session.UserID,
		session.Token,
		session.IPAddress,
		session.UserAgent,
		time.Now(),
		session.ExpiresAt,
		time.Now(),
		true,
	).Scan(&session.ID)
	return err
}

// GetSessionByToken obtiene una sesión por token
func (r *userRepository) GetSessionByToken(token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, ip_address, user_agent, created_at, expires_at, last_activity, active
		FROM sessions
		WHERE token = $1 AND active = true
	`
	session := &models.Session{}
	err := r.db.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.IPAddress,
		&session.UserAgent,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.LastActivity,
		&session.Active,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("sesión no encontrada")
		}
		return nil, err
	}
	return session, nil
}

// UpdateSessionActivity actualiza la última actividad de una sesión
func (r *userRepository) UpdateSessionActivity(sessionID int64) error {
	query := `UPDATE sessions SET last_activity = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), sessionID)
	return err
}

// DeactivateSession desactiva una sesión
func (r *userRepository) DeactivateSession(sessionID int64) error {
	query := `UPDATE sessions SET active = false WHERE id = $1`
	_, err := r.db.Exec(query, sessionID)
	return err
}

// GetUserSessions obtiene todas las sesiones de un usuario
func (r *userRepository) GetUserSessions(userID int64) ([]*models.Session, error) {
	query := `
		SELECT id, user_id, token, ip_address, user_agent, created_at, expires_at, last_activity, active
		FROM sessions
		WHERE user_id = $1 AND active = true
		ORDER BY last_activity DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Token,
			&session.IPAddress,
			&session.UserAgent,
			&session.CreatedAt,
			&session.ExpiresAt,
			&session.LastActivity,
			&session.Active,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

// CleanExpiredSessions limpia las sesiones expiradas
func (r *userRepository) CleanExpiredSessions() error {
	query := `UPDATE sessions SET active = false WHERE expires_at < $1 AND active = true`
	_, err := r.db.Exec(query, time.Now())
	return err
}
