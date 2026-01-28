package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"learning-go/domain"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running")
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Fprintf(w, "Error decoding JSON: %s", err)
		return
	}

	fmt.Fprintf(w, "User ID: %s\n", id)
	fmt.Fprintf(w, "User Name: %s\n", user.Name)
	fmt.Fprintf(w, "User Email: %s\n", user.Email)
	fmt.Fprintf(w, "User Age: %d\n", user.Age)
}
