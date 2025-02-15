package router

import (
	controllers "project/controllers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	controller := controllers.NewUserController()

	r.HandleFunc("/users", controller.GetAllUsersHandler).Methods("GET")

	return r
}
