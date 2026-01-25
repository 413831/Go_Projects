package http

import (
	"encoding/json"
	"kiosco/internal/application"
	"kiosco/internal/domain"
	"net/http"
)

// ItemHandler maneja las peticiones HTTP relacionadas con items
type ItemHandler struct {
	service *application.ItemService
}

// NewItemHandler crea una nueva instancia del handler
func NewItemHandler(service *application.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

// GetAllItems maneja GET /api/items
func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllItems(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, items)
}

// GetItemByID maneja GET /api/items/{id}
func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// Intentar obtener desde la URL path si está disponible
		id = r.PathValue("id")
	}
	
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID es requerido")
		return
	}
	
	item, err := h.service.GetItemByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrItemNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, item)
}

// CreateItem maneja POST /api/items
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item domain.Item
	
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error al decodificar el body: "+err.Error())
		return
	}
	
	createdItem, err := h.service.CreateItem(r.Context(), &item)
	if err != nil {
		if err == domain.ErrInvalidItemName || 
		   err == domain.ErrInvalidItemPrice || 
		   err == domain.ErrInvalidItemStock {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, createdItem)
}

// UpdateItem maneja PUT /api/items/{id}
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID es requerido")
		return
	}
	
	var item domain.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error al decodificar el body: "+err.Error())
		return
	}
	
	updatedItem, err := h.service.UpdateItem(r.Context(), id, &item)
	if err != nil {
		if err == domain.ErrItemNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == domain.ErrInvalidItemName || 
		   err == domain.ErrInvalidItemPrice || 
		   err == domain.ErrInvalidItemStock {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, updatedItem)
}

// DeleteItem maneja DELETE /api/items/{id}
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID es requerido")
		return
	}
	
	err := h.service.DeleteItem(r.Context(), id)
	if err != nil {
		if err == domain.ErrItemNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusNoContent, nil)
}

// respondWithJSON envía una respuesta JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError envía una respuesta de error JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
