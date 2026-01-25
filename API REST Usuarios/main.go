package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-rest-usuarios/config"
	"api-rest-usuarios/controllers"
	"api-rest-usuarios/repositories"
	"api-rest-usuarios/router"
	"api-rest-usuarios/services"
	"api-rest-usuarios/utils"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Inicializar logger
	logger := utils.GetLogger()
	defer logger.Close()

	logger.Info("Iniciando aplicación de gestión de usuarios...")

	// Conectar a la base de datos
	// Nota: El usuario debe configurar la conexión según su entorno
	// Por ahora comentamos la conexión real para que el usuario la configure
	/*
		db, err := database.Connect(cfg)
		if err != nil {
			logger.Error("Error al conectar a la base de datos: " + err.Error())
			log.Fatal(err)
		}
		defer db.Close()
		logger.Info("Conexión a base de datos establecida")
	*/

	// Para pruebas sin base de datos, usar mock repository
	// En producción, descomentar las líneas anteriores y usar:
	// userRepo := repositories.NewUserRepository(db)
	var userRepo repositories.UserRepository = repositories.NewMockUserRepository()
	logger.Info("Usando repositorio mock para pruebas (sin base de datos)")

	// Inicializar servicios de utilidad
	encryptionService, err := utils.NewEncryptionService(cfg.Security.EncryptionKey)
	if err != nil {
		logger.Error("Error al inicializar servicio de encriptación: " + err.Error())
		log.Fatal(err)
	}

	// Inicializar servicios
	userService := services.NewUserService(userRepo, encryptionService, logger, cfg.Security)
	sessionService := services.NewSessionService(userRepo, logger)

	// Inicializar controladores
	userController := controllers.NewUserController(userService, sessionService, logger)

	// Configurar router
	r := router.SetupRouter(userController)

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		logger.Info(fmt.Sprintf("Servidor iniciado en http://%s:%s", cfg.Server.Host, cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error al iniciar servidor: " + err.Error())
			log.Fatal(err)
		}
	}()

	// Esperar señal de interrupción para apagado graceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Apagando servidor...")

	// Detener servicio de sesiones
	sessionService.Stop()

	// Apagar servidor con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error al apagar servidor: " + err.Error())
		log.Fatal(err)
	}

	logger.Info("Servidor apagado correctamente")
}
