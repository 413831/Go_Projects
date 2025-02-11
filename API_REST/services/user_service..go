package services

import respositories "project/repositories"

type UserService struct {
	repo respositories.UserRepository
}

func NewUserService() *UserService {
	repo := &respositories.MockUserRepository{}

	return &UserService{repo: repo}
}

func (s *UserService) GetUsers() []respositories.User {
	return s.repo.GetAllUsers()
}
