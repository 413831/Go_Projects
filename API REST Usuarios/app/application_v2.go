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
	"api-rest-usuarios/container"
	"api-rest-usuarios/controllers"
	"api-rest-usuarios/repositories"
	"api-rest-usuarios/router"
	"api-rest-usuarios/services"
	"api-rest-usuarios/utils"
)

// ApplicationV2 implementa la aplicación con patrones de diseño avanzados
type ApplicationV2 struct {
	config            *config.Config
	logger            *utils.Logger
	diContainer       container.Container
	eventPublisher    services.EventPublisher
	server            *http.Server
}

// NewApplicationV2 crea una nueva instancia de la aplicación V2
func NewApplicationV2(cfg *config.Config) (*ApplicationV2, error) {
	logger := utils.GetLogger()
	logger.Info("Iniciando aplicación V2 con patrones de diseño...")

	// Crear contenedor de inyección de dependencias
	diContainer := container.NewDIContainer()

	// Crear publicador de eventos
	eventPublisher := services.NewEventPublisher()

	// Registrar observadores
	loggerObserver := services.NewLoggerObserver("UserService")
	auditObserver := services.NewAuditObserver()

	eventPublisher.Subscribe("user.created", loggerObserver)
	eventPublisher.Subscribe("user.updated", loggerObserver)
	eventPublisher.Subscribe("user.deleted", loggerObserver)
	eventPublisher.Subscribe("session.created", loggerObserver)
	eventPublisher.Subscribe("error.occurred", loggerObserver)

	eventPublisher.Subscribe("user.created", auditObserver)
	eventPublisher.Subscribe("user.updated", auditObserver)
	eventPublisher.Subscribe("user.deleted", auditObserver)

	// Crear fábrica de repositorios
	repoFactory := services.NewMockRepositoryFactory()

	// Registrar repositorios en el contenedor
	diContainer.RegisterSingleton("userRepo", func() (interface{}, error) {
		return repoFactory.CreateUserRepository(), nil
	})
	diContainer.RegisterSingleton("roleRepo", func() (interface{}, error) {
		return repoFactory.CreateRoleRepository(), nil
	})
	diContainer.RegisterSingleton("piiRepo", func() (interface{}, error) {
		return repoFactory.CreatePIIRepository(), nil
	})
	diContainer.RegisterSingleton("sessionRepo", func() (interface{}, error) {
		return repoFactory.CreateSessionRepository(), nil
	})

	// Registrar estrategias
	diContainer.RegisterSingleton("passwordHasher", func() (interface{}, error) {
		return services.NewBCryptPasswordHasher(cfg.Security.GetBCryptCost()), nil
	})
	diContainer.RegisterSingleton("piiEncryptor", func() (interface{}, error) {
		return services.NewAESEncryptor(cfg.Security.EncryptionKey)
	})
	diContainer.RegisterSingleton("validator", func() (interface{}, error) {
		return services.NewUserValidator(), nil
	})

	// Registrar servicios principales
	diContainer.RegisterSingleton("eventPublisher", func() (interface{}, error) {
		return eventPublisher, nil
	})

	diContainer.RegisterSingleton("userService", func() (interface{}, error) {
		userRepo, err := diContainer.Get("userRepo")
		if err != nil {
			return nil, err
		}
		roleRepo, err := diContainer.Get("roleRepo")
		if err != nil {
			return nil, err
		}
		piiRepo, err := diContainer.Get("piiRepo")
		if err != nil {
			return nil, err
		}
		hasher, err := diContainer.Get("passwordHasher")
		if err != nil {
			return nil, err
		}
		encryptor, err := diContainer.Get("piiEncryptor")
		if err != nil {
			return nil, err
		}
		validator, err := diContainer.Get("validator")
		if err != nil {
			return nil, err
		}
		publisher, err := diContainer.Get("eventPublisher")
		if err != nil {
			return nil, err
		}

		return services.NewUserServiceV2(
			userRepo.(repositories.UserRepository),
			roleRepo.(repositories.RoleRepository),
			piiRepo.(repositories.PIIRepository),
			hasher.(services.PasswordHasher),
			encryptor.(services.PIIEncryptor),
			validator.(services.Validator),
			publisher.(services.EventPublisher),
		), nil
	})

	diContainer.RegisterSingleton("sessionService", func() (interface{}, error) {
		sessionRepo, err := diContainer.Get("sessionRepo")
		if err != nil {
			return nil, err
		}
		userRepo, err := diContainer.Get("userRepo")
		if err != nil {
			return nil, err
		}
		publisher, err := diContainer.Get("eventPublisher")
		if err != nil {
			return nil, err
		}

		return services.NewSessionServiceV2(
			sessionRepo.(repositories.SessionRepository),
			userRepo.(repositories.UserRepository),
			publisher.(services.EventPublisher),
		), nil
	})

	diContainer.RegisterSingleton("roleService", func() (interface{}, error) {
		roleRepo, err := diContainer.Get("roleRepo")
		if err != nil {
			return nil, err
		}
		publisher, err := diContainer.Get("eventPublisher")
		if err != nil {
			return nil, err
		}

		return services.NewRoleService(
			roleRepo.(repositories.RoleRepository),
			publisher.(services.EventPublisher),
		), nil
	})

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &ApplicationV2{
		config:         cfg,
		logger:         logger,
		diContainer:    diContainer,
		eventPublisher: eventPublisher,
		server:         server,
	}, nil
}

// Start inicia la aplicación
func (app *ApplicationV2) Start() error {
	// Obtener servicios del contenedor
	userService, err := app.diContainer.Get("userService")
	if err != nil {
		return fmt.Errorf("error obteniendo userService: %w", err)
	}

	sessionService, err := app.diContainer.Get("sessionService")
	if err != nil {
		return fmt.Errorf("error obteniendo sessionService: %w", err)
	}

	// Crear controladores
	userController := controllers.NewUserControllerV2(
		userService.(services.UserServiceInterface),
		sessionService.(services.SessionServiceInterface),
		app.logger,
	)

	// Configurar router
	handler := router.SetupRouter(userController)
	app.server.Handler = handler

	// Iniciar servidor en goroutine
	go func() {
		app.logger.Info(fmt.Sprintf("Servidor V2 iniciado en http://%s:%s", app.config.Server.Host, app.config.Server.Port))
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Error al iniciar servidor: " + err.Error())
			log.Fatal(err)
		}
	}()

	return nil
}

// WaitForShutdown espera la señal de apagado
func (app *ApplicationV2) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.logger.Info("Apagando servidor V2...")

	// Detener servicios
	if sessionService, err := app.diContainer.Get("sessionService"); err == nil {
		if ss, ok := sessionService.(services.SessionServiceInterface); ok {
			ss.Stop()
		}
	}

	// Apagar servidor con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error("Error al apagar servidor: " + err.Error())
		log.Fatal(err)
	}

	app.logger.Info("Servidor V2 apagado correctamente")
}

// Close cierra recursos
func (app *ApplicationV2) Close() {
	if app.logger != nil {
		_ = app.logger.Close()
	}
}

// GetContainer retorna el contenedor de inyección de dependencias
func (app *ApplicationV2) GetContainer() container.Container {
	return app.diContainer
}

// GetEventPublisher retorna el publicador de eventos
func (app *ApplicationV2) GetEventPublisher() services.EventPublisher {
	return app.eventPublisher
}
