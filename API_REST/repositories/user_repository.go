package respositories

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserRepository interface {
	GetAllUsers() []User
}

type MockUserRepository struct{}

func (m *MockUserRepository) GetAllUsers() []User {
	return []User{
		{
			ID:   1,
			Name: "Juan Perez",
		},
		{
			ID:   2,
			Name: "Maria Lopez",
		},
	}
}
