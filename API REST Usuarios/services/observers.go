package services

import (
	"fmt"
	"sync"
	"time"
)

// Event define la interfaz base para eventos
type Event interface {
	Type() string
	Data() interface{}
	Timestamp() time.Time
}

// Observer define la interfaz para observadores
type Observer interface {
	Notify(event Event)
}

// EventPublisher define la interfaz para publicar eventos
type EventPublisher interface {
	Subscribe(eventType string, observer Observer)
	Unsubscribe(eventType string, observer Observer)
	Publish(event Event)
}

// UserEvent representa eventos relacionados con usuarios
type UserEvent struct {
	eventType string
	data      interface{}
	timestamp time.Time
}

// NewUserEvent crea un nuevo evento de usuario
func NewUserEvent(eventType string, data interface{}) Event {
	return &UserEvent{
		eventType: eventType,
		data:      data,
		timestamp: time.Now(),
	}
}

func (e *UserEvent) Type() string {
	return e.eventType
}

func (e *UserEvent) Data() interface{} {
	return e.data
}

func (e *UserEvent) Timestamp() time.Time {
	return e.timestamp
}

// EventPublisherImpl implementa EventPublisher
type EventPublisherImpl struct {
	observers map[string][]Observer
	mu        sync.RWMutex
}

// NewEventPublisher crea un nuevo publicador de eventos
func NewEventPublisher() EventPublisher {
	return &EventPublisherImpl{
		observers: make(map[string][]Observer),
	}
}

func (ep *EventPublisherImpl) Subscribe(eventType string, observer Observer) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	if ep.observers[eventType] == nil {
		ep.observers[eventType] = make([]Observer, 0)
	}
	ep.observers[eventType] = append(ep.observers[eventType], observer)
}

func (ep *EventPublisherImpl) Unsubscribe(eventType string, observer Observer) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	observers := ep.observers[eventType]
	for i, obs := range observers {
		if obs == observer {
			ep.observers[eventType] = append(observers[:i], observers[i+1:]...)
			break
		}
	}
}

func (ep *EventPublisherImpl) Publish(event Event) {
	ep.mu.RLock()
	defer ep.mu.RUnlock()

	observers := ep.observers[event.Type()]
	for _, observer := range observers {
		go observer.Notify(event) // Notificar asíncronamente
	}
}

// LoggerObserver implementa Observer para logging
type LoggerObserver struct {
	name string
}

// NewLoggerObserver crea un nuevo observador de logging
func NewLoggerObserver(name string) Observer {
	return &LoggerObserver{name: name}
}

func (lo *LoggerObserver) Notify(event Event) {
	switch event.Type() {
	case "user.created":
		if user, ok := event.Data().(*UserData); ok {
			fmt.Printf("[%s] Usuario creado: %s (ID: %d) at %v\n", 
				lo.name, user.Username, user.UserID, event.Timestamp())
		}
	case "user.updated":
		if user, ok := event.Data().(*UserData); ok {
			fmt.Printf("[%s] Usuario actualizado: %s (ID: %d) at %v\n", 
				lo.name, user.Username, user.UserID, event.Timestamp())
		}
	case "user.deleted":
		if user, ok := event.Data().(*UserData); ok {
			fmt.Printf("[%s] Usuario eliminado: %s (ID: %d) at %v\n", 
				lo.name, user.Username, user.UserID, event.Timestamp())
		}
	case "session.created":
		if session, ok := event.Data().(*SessionData); ok {
			fmt.Printf("[%s] Sesión creada: %s (UserID: %d) at %v\n", 
				lo.name, session.Token, session.UserID, event.Timestamp())
		}
	case "error.occurred":
		if err, ok := event.Data().(*ErrorData); ok {
			fmt.Printf("[%s] Error: %s at %v\n", 
				lo.name, err.Message, event.Timestamp())
		}
	}
}

// AuditObserver implementa Observer para auditoría
type AuditObserver struct {
	auditLog []AuditEntry
	mu       sync.Mutex
}

// AuditEntry representa una entrada de auditoría
type AuditEntry struct {
	Timestamp time.Time
	EventType string
	UserID    int64
	Action    string
	Details   string
}

// NewAuditObserver crea un nuevo observador de auditoría
func NewAuditObserver() Observer {
	return &AuditObserver{
		auditLog: make([]AuditEntry, 0),
	}
}

func (ao *AuditObserver) Notify(event Event) {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	entry := AuditEntry{
		Timestamp: event.Timestamp(),
		EventType: event.Type(),
	}

	switch event.Type() {
	case "user.created":
		if user, ok := event.Data().(*UserData); ok {
			entry.UserID = user.UserID
			entry.Action = "CREATE_USER"
			entry.Details = fmt.Sprintf("Username: %s, Email: %s", user.Username, user.Email)
		}
	case "user.updated":
		if user, ok := event.Data().(*UserData); ok {
			entry.UserID = user.UserID
			entry.Action = "UPDATE_USER"
			entry.Details = fmt.Sprintf("Username: %s", user.Username)
		}
	case "user.deleted":
		if user, ok := event.Data().(*UserData); ok {
			entry.UserID = user.UserID
			entry.Action = "DELETE_USER"
			entry.Details = fmt.Sprintf("Username: %s", user.Username)
		}
	}

	ao.auditLog = append(ao.auditLog, entry)
}

// GetAuditLog retorna el log de auditoría
func (ao *AuditObserver) GetAuditLog() []AuditEntry {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	// Retornar copia para evitar modificaciones externas
	log := make([]AuditEntry, len(ao.auditLog))
	copy(log, ao.auditLog)
	return log
}

// UserData representa datos de usuario para eventos
type UserData struct {
	UserID   int64
	Username string
	Email    string
}

// SessionData representa datos de sesión para eventos
type SessionData struct {
	SessionID int64
	UserID    int64
	Token     string
	IPAddress string
}

// ErrorData representa datos de error para eventos
type ErrorData struct {
	UserID  int64
	Message string
	Context string
}
