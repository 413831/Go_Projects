package services

import (
	models "project/models"
	repositories "project/repositories"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService() *UserService {
	repo := &repositories.MockUserRepository{}

	return &UserService{repo: repo}
}

func (s *UserService) GetUsers() []models.User {
	return s.repo.GetAllUsers()
}
