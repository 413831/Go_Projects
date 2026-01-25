package router

import (
	"api-rest-usuarios/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRouter configura todas las rutas de la aplicación
func SetupRouter(userController *controllers.UserController) *mux.Router {
	r := mux.NewRouter()

	// Middleware para CORS y logging
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	// Rutas de usuarios
	api := r.PathPrefix("/api/v1").Subrouter()

	// CRUD de usuarios
	api.HandleFunc("/users", userController.CreateUser).Methods("POST")
	api.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", userController.GetUser).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", userController.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id:[0-9]+}", userController.DeleteUser).Methods("DELETE")

	// Gestión de roles
	api.HandleFunc("/users/{id:[0-9]+}/roles", userController.GrantRole).Methods("POST")
	api.HandleFunc("/users/{id:[0-9]+}/roles/{role_id:[0-9]+}", userController.RevokeRole).Methods("DELETE")

	// Sesiones
	api.HandleFunc("/users/{id:[0-9]+}/sessions", userController.GetUserSessions).Methods("GET")

	// Health check
	r.HandleFunc("/health", healthCheck).Methods("GET")

	return r
}

// loggingMiddleware registra las peticiones HTTP
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log básico de peticiones
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware maneja CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// healthCheck endpoint para verificar el estado del servicio
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
