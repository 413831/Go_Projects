package services

import (
	"api-rest-usuarios/config"
	"api-rest-usuarios/repositories"
)

// ServiceFactory define la interfaz para crear servicios
type ServiceFactory interface {
	CreateUserService() UserServiceInterface
	CreateSessionService() SessionServiceInterface
	CreateRoleService() RoleServiceInterface
}

// ServiceFactoryImpl implementa ServiceFactory
type ServiceFactoryImpl struct {
	userRepo    repositories.UserRepository
	roleRepo    repositories.RoleRepository
	piiRepo     repositories.PIIRepository
	sessionRepo repositories.SessionRepository
	config      *config.Config
	eventPublisher EventPublisher
}

// NewServiceFactory crea una nueva fábrica de servicios
func NewServiceFactory(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	piiRepo repositories.PIIRepository,
	sessionRepo repositories.SessionRepository,
	config *config.Config,
	eventPublisher EventPublisher,
) ServiceFactory {
	return &ServiceFactoryImpl{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		piiRepo:        piiRepo,
		sessionRepo:    sessionRepo,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

func (f *ServiceFactoryImpl) CreateUserService() UserServiceInterface {
	// Crear estrategias
	hasher := NewBCryptPasswordHasher(f.config.Security.GetBCryptCost())
	encryptor, _ := NewAESEncryptor(f.config.Security.EncryptionKey)
	validator := NewUserValidator()

	// Crear servicio con todas las dependencias
	return NewUserServiceV2(
		f.userRepo,
		f.roleRepo,
		f.piiRepo,
		hasher,
		encryptor,
		validator,
		f.eventPublisher,
	)
}

func (f *ServiceFactoryImpl) CreateSessionService() SessionServiceInterface {
	return NewSessionServiceV2(f.sessionRepo, f.userRepo, f.eventPublisher)
}

func (f *ServiceFactoryImpl) CreateRoleService() RoleServiceInterface {
	return NewRoleService(f.roleRepo, f.eventPublisher)
}

// RepositoryFactory define la interfaz para crear repositorios
type RepositoryFactory interface {
	CreateUserRepository() repositories.UserRepository
	CreateRoleRepository() repositories.RoleRepository
	CreatePIIRepository() repositories.PIIRepository
	CreateSessionRepository() repositories.SessionRepository
	CreateCompositeRepository() repositories.CompositeRepository
}

// MockRepositoryFactory implementa RepositoryFactory para mocks
type MockRepositoryFactory struct{}

// NewMockRepositoryFactory crea una nueva fábrica de repositorios mock
func NewMockRepositoryFactory() RepositoryFactory {
	return &MockRepositoryFactory{}
}

func (f *MockRepositoryFactory) CreateUserRepository() repositories.UserRepository {
	return repositories.NewMockUserRepositoryImpl()
}

func (f *MockRepositoryFactory) CreateRoleRepository() repositories.RoleRepository {
	return repositories.NewMockRoleRepositoryImpl()
}

func (f *MockRepositoryFactory) CreatePIIRepository() repositories.PIIRepository {
	return repositories.NewMockPIIRepositoryImpl()
}

func (f *MockRepositoryFactory) CreateSessionRepository() repositories.SessionRepository {
	return repositories.NewMockSessionRepositoryImpl()
}

func (f *MockRepositoryFactory) CreateCompositeRepository() repositories.CompositeRepository {
	return repositories.NewMockCompositeRepository()
}

// PostgresRepositoryFactory implementa RepositoryFactory para PostgreSQL
type PostgresRepositoryFactory struct {
	db interface{} // *sql.DB en una implementación real
}

// NewPostgresRepositoryFactory crea una nueva fábrica de repositorios PostgreSQL
func NewPostgresRepositoryFactory(db interface{}) RepositoryFactory {
	return &PostgresRepositoryFactory{db: db}
}

func (f *PostgresRepositoryFactory) CreateUserRepository() repositories.UserRepository {
	return repositories.NewPostgresUserRepository(f.db) // Implementación pendiente
}

func (f *PostgresRepositoryFactory) CreateRoleRepository() repositories.RoleRepository {
	return repositories.NewPostgresRoleRepository(f.db) // Implementación pendiente
}

func (f *PostgresRepositoryFactory) CreatePIIRepository() repositories.PIIRepository {
	return repositories.NewPostgresPIIRepository(f.db) // Implementación pendiente
}

func (f *PostgresRepositoryFactory) CreateSessionRepository() repositories.SessionRepository {
	return repositories.NewPostgresSessionRepository(f.db) // Implementación pendiente
}

func (f *PostgresRepositoryFactory) CreateCompositeRepository() repositories.CompositeRepository {
	return repositories.NewPostgresCompositeRepository(f.db)
}
