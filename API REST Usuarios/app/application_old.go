//go:build ignore

package app

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

type Application struct {
	config            *config.Config
	logger            *utils.Logger
	userRepo          repositories.UserRepository
	encryptionService *utils.EncryptionService
	userService       *services.UserService
	sessionService    *services.SessionService
	userController    *controllers.UserController
	handler           http.Handler
	server            *http.Server
}

func NewApplication(cfg *config.Config) (*Application, error) {
	logger := utils.GetLogger()
	logger.Info("Iniciando aplicación de gestión de usuarios...")

	// Para pruebas sin base de datos, usar mock repository
	// En producción, descomentar las líneas anteriores y usar:
	// userRepo := repositories.NewUserRepository(db)
	var userRepo repositories.UserRepository = repositories.NewMockUserRepository()
	logger.Info("Usando repositorio mock para pruebas (sin base de datos)")

	encryptionService, err := utils.NewEncryptionService(cfg.Security.EncryptionKey)
	if err != nil {
		logger.Error("Error al inicializar servicio de encriptación: " + err.Error())
		return nil, err
	}

	userService := services.NewUserService(userRepo, encryptionService, logger, cfg.Security)
	sessionService := services.NewSessionService(userRepo, logger)
	userController := controllers.NewUserController(userService, sessionService, logger)

	handler := router.SetupRouter(userController)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Application{
		config:            cfg,
		logger:            logger,
		userRepo:          userRepo,
		encryptionService: encryptionService,
		userService:       userService,
		sessionService:    sessionService,
		userController:    userController,
		handler:           handler,
		server:            server,
	}, nil
}

func (app *Application) Start() {
	go func() {
		app.logger.Info(fmt.Sprintf("Servidor iniciado en http://%s:%s", app.config.Server.Host, app.config.Server.Port))
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Error al iniciar servidor: " + err.Error())
			log.Fatal(err)
		}
	}()
}

func (app *Application) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.logger.Info("Apagando servidor...")
	app.sessionService.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error("Error al apagar servidor: " + err.Error())
		log.Fatal(err)
	}

	app.logger.Info("Servidor apagado correctamente")
}

func (app *Application) Close() {
	if app.logger != nil {
		_ = app.logger.Close()
	}
}
