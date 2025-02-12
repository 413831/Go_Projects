package respositories

import models "project/models"

type UserRepository interface {
	GetAllUsers() []models.User
}

type MockUserRepository struct{}

func (m *MockUserRepository) GetAllUsers() []models.User {
	return []models.User{
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
