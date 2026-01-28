package internal

import (
	"errors"
	"learning-go/domain"
)

type UserRepository interface {
	GetByID(id string) (*domain.User, error)
}

type MockRepository struct {
	users map[string]*domain.User
}

func NewMockRepo() *MockRepository {
	return &MockRepository{
		users: map[string]*domain.User{
			"1": {
				Id:    1,
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Age:   30,
			},
			"2": {
				Id:    2,
				Name:  "Jane Doe",
				Email: "jane.doe@example.com",
				Age:   25,
			},
		}}
}

func (r *MockRepository) GetByID(id string) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("User not found")
	}

	return user, nil
}
