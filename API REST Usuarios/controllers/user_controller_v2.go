package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-rest-usuarios/models"
	"api-rest-usuarios/services"
	"api-rest-usuarios/utils"

	"github.com/gorilla/mux"
)

// UserControllerV2 implementa el controlador con patrones de diseño
type UserControllerV2 struct {
	userService    services.UserServiceInterface
	sessionService services.SessionServiceInterface
	logger         *utils.Logger
}

// NewUserControllerV2 crea un nuevo controlador de usuarios V2
func NewUserControllerV2(
	userService services.UserServiceInterface,
	sessionService services.SessionServiceInterface,
	logger *utils.Logger,
) *UserControllerV2 {
	return &UserControllerV2{
		userService:    userService,
		sessionService: sessionService,
		logger:         logger,
	}
}

func (uc *UserControllerV2) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "Solicitud inválida")
		return
	}

	user, err := uc.userService.Create(&req)
	if err != nil {
		uc.logger.Error("Error al crear usuario: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Usar Builder Pattern para construir respuesta
	response := services.NewUserResponseBuilder(user).
		WithRoles(user.Roles).
		WithPII(user.PII).
		Build()

	uc.respondWithJSON(w, http.StatusCreated, response)
}

func (uc *UserControllerV2) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	user, err := uc.userService.GetByID(id)
	if err != nil {
		uc.respondWithError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	// Usar Builder Pattern para construir respuesta
	response := services.NewUserResponseBuilder(user).
		WithRoles(user.Roles).
		WithPII(user.PII).
		Build()

	uc.respondWithJSON(w, http.StatusOK, response)
}

func (uc *UserControllerV2) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	users, total, err := uc.userService.GetAll(page, pageSize)
	if err != nil {
		uc.logger.Error("Error al obtener usuarios: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, "Error al obtener usuarios")
		return
	}

	// Construir respuestas para cada usuario
	var userResponses []interface{}
	for _, user := range users {
		response := services.NewUserResponseBuilder(user).
			WithRoles(user.Roles).
			Build()
		userResponses = append(userResponses, response)
	}

	// Usar Builder Pattern para respuesta paginada
	response := services.NewPaginatedResponseBuilder(userResponses, page, pageSize, total).Build()
	uc.respondWithJSON(w, http.StatusOK, response)
}

func (uc *UserControllerV2) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "Solicitud inválida")
		return
	}

	user, err := uc.userService.Update(id, &req)
	if err != nil {
		uc.logger.Error("Error al actualizar usuario: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Usar Builder Pattern para construir respuesta
	response := services.NewUserResponseBuilder(user).
		WithRoles(user.Roles).
		WithPII(user.PII).
		Build()

	uc.respondWithJSON(w, http.StatusOK, response)
}

func (uc *UserControllerV2) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	err = uc.userService.Delete(id)
	if err != nil {
		uc.logger.Error("Error al eliminar usuario: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, "Error al eliminar usuario")
		return
	}

	uc.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Usuario eliminado correctamente"})
}

func (uc *UserControllerV2) GrantRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req models.GrantRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "Solicitud inválida")
		return
	}

	grantedBy := userID
	err = uc.userService.GrantRole(userID, req.RoleID, grantedBy)
	if err != nil {
		uc.logger.Error("Error al otorgar rol: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, "Error al otorgar rol")
		return
	}

	uc.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Rol otorgado correctamente"})
}

func (uc *UserControllerV2) RevokeRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	roleID, err := strconv.ParseInt(vars["role_id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "Role ID inválido")
		return
	}

	err = uc.userService.RevokeRole(userID, roleID)
	if err != nil {
		uc.logger.Error("Error al revocar rol: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, "Error al revocar rol")
		return
	}

	uc.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Rol revocado correctamente"})
}

func (uc *UserControllerV2) GetUserSessions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		uc.respondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	sessions, err := uc.sessionService.GetUserSessions(userID)
	if err != nil {
		uc.logger.Error("Error al obtener sesiones: " + err.Error())
		uc.respondWithError(w, http.StatusInternalServerError, "Error al obtener sesiones")
		return
	}

	uc.respondWithJSON(w, http.StatusOK, sessions)
}

func (uc *UserControllerV2) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (uc *UserControllerV2) respondWithError(w http.ResponseWriter, code int, message string) {
	uc.respondWithJSON(w, code, map[string]string{"error": message})
}
