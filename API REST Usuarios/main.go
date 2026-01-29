package main

import (
	"log"

	"api-rest-usuarios/app"
	"api-rest-usuarios/config"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Crear aplicación con patrones de diseño
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	// Iniciar aplicación
	if err := application.Start(); err != nil {
		log.Fatal(err)
	}

	// Esperar señal de apagado
	application.WaitForShutdown()
}
