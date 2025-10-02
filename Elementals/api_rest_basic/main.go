package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Country representa la estructura de datos para un país
// Contiene el nombre del país y su idioma principal
type Country struct {
	Name     string // Nombre del país
	Language string // Idioma principal del país
}

// countries es un slice global que almacena la lista de países
var countries []*Country

// main es la función principal que inicia un servidor HTTP básico
// Configura el manejo de señales para un cierre graceful del servidor
func main() {
	// Crea un contexto de fondo para el shutdown del servidor
	ctx := context.Background()

	// Canal para recibir señales del sistema operativo
	sc := make(chan os.Signal, 1)

	// Configura el canal para recibir señales de interrupción y terminación
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	// Configura las rutas del servidor
	initRoutes()

	// Crea una nueva instancia del servidor en el puerto 8080
	srv := &http.Server{
		Addr: ":8080",
	}

	// Inicia el servidor en una goroutine separada
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	log.Println("server started")

	// Espera a recibir una señal para cerrar el servidor
	<-sc

	// Realiza un shutdown graceful del servidor
	srv.Shutdown(ctx)
	log.Println("server stopped")
}

// initRoutes configura todas las rutas del servidor HTTP
// Registra los handlers para las diferentes rutas y métodos HTTP
func initRoutes() {
	// Ruta raíz que maneja solo peticiones GET
	http.HandleFunc("/", index)

	// Ruta /countries que maneja tanto GET como POST
	http.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCountries(w, r)
		case http.MethodPost:
			addCountry(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method not allowed")
			return
		}
	})
}

// index maneja las peticiones GET a la ruta raíz "/"
// Responde con un mensaje de bienvenida al visitante
func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed")
		return
	}
	fmt.Fprint(w, "Hello there %s", "visitor")
}

// getCountries maneja las peticiones GET a "/countries"
// Retorna la lista completa de países en formato JSON
func getCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

// addCountry maneja las peticiones POST a "/countries"
// Decodifica el JSON del cuerpo de la petición y agrega un nuevo país a la lista
func addCountry(w http.ResponseWriter, r *http.Request) {
	country := &Country{}

	err := json.NewDecoder(r.Body).Decode(country)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}
	countries = append(countries, country)
	fmt.Fprintf(w, "country was added")
}
