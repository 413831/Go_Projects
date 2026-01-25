package utils

import (
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel representa el nivel de log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger maneja el logging de la aplicación
type Logger struct {
	file   *os.File
	logger *log.Logger
	mu     sync.Mutex
	level  LogLevel
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger retorna una instancia singleton del logger
func GetLogger() *Logger {
	once.Do(func() {
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Error al abrir archivo de log:", err)
		}
		instance = &Logger{
			file:   file,
			logger: log.New(file, "", log.LstdFlags),
			level:  INFO,
		}
	})
	return instance
}

// Log registra un mensaje con el nivel especificado
func (l *Logger) Log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	levelStr := []string{"DEBUG", "INFO", "WARN", "ERROR"}[level]
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.logger.Printf("[%s] [%s] %s", timestamp, levelStr, message)
}

// Debug registra un mensaje de debug
func (l *Logger) Debug(message string) {
	l.Log(DEBUG, message)
}

// Info registra un mensaje informativo
func (l *Logger) Info(message string) {
	l.Log(INFO, message)
}

// Warn registra un mensaje de advertencia
func (l *Logger) Warn(message string) {
	l.Log(WARN, message)
}

// Error registra un mensaje de error
func (l *Logger) Error(message string) {
	l.Log(ERROR, message)
}

// Close cierra el archivo de log
func (l *Logger) Close() error {
	return l.file.Close()
}
