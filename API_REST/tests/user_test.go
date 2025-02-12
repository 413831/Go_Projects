package tests

import (
	"project/services"
	"testing"
)

func TestGetUsers(t *testing.T) {
	service := services.NewUserService()
	users := service.GetUsers()

	if len(users) == 0 {
		t.Errorf("Expected users, got none")
	}
}
