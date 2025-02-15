package controllers

import (
	"net/http"
	"project/services"
	"project/utils"
)

type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	service := services.NewUserService()

	return &UserController{service: service}
}

func (c *UserController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := c.service.GetUsers()

	utils.JSONResponse(w, users, http.StatusOK)
}
