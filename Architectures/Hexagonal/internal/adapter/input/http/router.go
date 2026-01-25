package http

import (
	"kiosco/internal/application"
	"net/http"

	"github.com/gorilla/mux"
)

// Router configura las rutas de la API
type Router struct {
	itemHandler *ItemHandler
}

// NewRouter crea una nueva instancia del router
func NewRouter(itemService *application.ItemService) *Router {
	return &Router{
		itemHandler: NewItemHandler(itemService),
	}
}

// SetupRoutes configura todas las rutas de la API
func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	
	// Rutas de items
	api := router.PathPrefix("/api").Subrouter()
	
	// GET /api/items - Obtener todos los items
	api.HandleFunc("/items", r.itemHandler.GetAllItems).Methods("GET")
	
	// GET /api/items/{id} - Obtener un item por ID
	api.HandleFunc("/items/{id}", r.itemHandler.GetItemByID).Methods("GET")
	
	// POST /api/items - Crear un nuevo item
	api.HandleFunc("/items", r.itemHandler.CreateItem).Methods("POST")
	
	// PUT /api/items/{id} - Actualizar un item
	api.HandleFunc("/items/{id}", r.itemHandler.UpdateItem).Methods("PUT")
	
	// DELETE /api/items/{id} - Eliminar un item
	api.HandleFunc("/items/{id}", r.itemHandler.DeleteItem).Methods("DELETE")
	
	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	return router
}
