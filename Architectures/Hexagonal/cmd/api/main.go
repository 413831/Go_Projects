package main

import (
	"fmt"
	httphandler "kiosco/internal/adapter/input/http"
	"kiosco/internal/adapter/output/api"
	"kiosco/internal/application"
	"kiosco/internal/config"
	"log"
	"net/http"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Inicializar adaptador de salida (repositorio HTTP)
	itemRepository := api.NewHTTPItemRepository(cfg.ExternalAPIURL)

	// Inicializar casos de uso (servicio de aplicación)
	itemService := application.NewItemService(itemRepository)

	// Inicializar adaptador de entrada (router HTTP)
	router := httphandler.NewRouter(itemService)
	muxRouter := router.SetupRoutes()

	// Iniciar servidor
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Servidor iniciado en http://localhost%s", serverAddr)
	log.Printf("API disponible en http://localhost%s/api", serverAddr)
	log.Printf("Conectado a API externa: %s", cfg.ExternalAPIURL)

	if err := http.ListenAndServe(serverAddr, muxRouter); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
