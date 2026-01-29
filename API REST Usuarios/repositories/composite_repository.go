package repositories

// CompositeRepositoryImpl implementa CompositeRepository combinando todos los repositorios
type CompositeRepositoryImpl struct {
	UserRepository
	RoleRepository
	PIIRepository
	SessionRepository
}

// NewMockCompositeRepository crea un repositorio compuesto mock
func NewMockCompositeRepository() CompositeRepository {
	return &CompositeRepositoryImpl{
		UserRepository:    NewMockUserRepositoryImpl(),
		RoleRepository:    NewMockRoleRepositoryImpl(),
		PIIRepository:     NewMockPIIRepositoryImpl(),
		SessionRepository: NewMockSessionRepositoryImpl(),
	}
}

// NewPostgresCompositeRepository crea un repositorio compuesto para PostgreSQL
func NewPostgresCompositeRepository(db interface{}) CompositeRepository {
	// Aquí iría la implementación para PostgreSQL
	// Por ahora retornamos el mock como placeholder
	return NewMockCompositeRepository()
}
