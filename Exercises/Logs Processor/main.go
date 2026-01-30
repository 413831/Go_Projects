package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type LogStats struct {
	Errors   int
	Warnings int
	Info     int
	mu       sync.Mutex
}

func (ls *LogStats) Increment(level string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	switch level {
	case "ERROR":
		ls.Errors++
	case "WARNING":
		ls.Warnings++
	case "INFO":
		ls.Info++
	}
}

func (ls *LogStats) Print() {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	total := ls.Errors + ls.Warnings + ls.Info

	fmt.Printf("\n Estadísticas:\n")
	fmt.Printf(" Errores: %d - %d%% \n", ls.Errors, ls.Errors*100/total)
	fmt.Printf(" Warning: %d - %d%% \n", ls.Warnings, ls.Warnings*100/total)
	fmt.Printf(" Info: %d - %d%% \n", ls.Info, ls.Info*100/total)
	fmt.Printf(" Total: %d", total)
}

func processLine(line string, stats *LogStats, alerts chan<- string) {
	line = strings.TrimSpace(line)

	if line == "" {
		return
	}

	if strings.Contains(line, "ERROR") {
		stats.Increment("ERROR")
		alerts <- fmt.Sprintf("ERROR : %s", line)
	} else if strings.Contains(line, "WARNING") {
		stats.Increment("WARNING")
		alerts <- fmt.Sprintf("WARNING : %s", line)
	} else if strings.Contains(line, "INFO") {
		stats.Increment("INFO")
	}
}

func main() {
	var (
		stats      = &LogStats{}
		alerts     = make(chan string, 100)
		logFile, _ = os.Create("app.log")
	)

	logFile.WriteString("2024-01-01 10:00:00 INFO Application started\n")
	logFile.WriteString("2024-01-01 10:00:05 INFO User logged in\n")
	logFile.WriteString("2024-01-01 10:00:10 WARNING Low memory\n")
	logFile.WriteString("2024-01-01 10:00:15 ERROR Database connection failed\n")
	logFile.WriteString("2024-01-01 10:00:20 INFO Request processed\n")
	logFile.WriteString("2024-01-01 10:00:25 ERROR Timeout exceeded\n")
	logFile.Close()

	// Leer archivo
	file, err := os.Open("app.log")
	if err != nil {
		fmt.Println("Error abriendo archivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup

	// Procesar cada línea concurrentemente
	for scanner.Scan() {
		wg.Add(1)
		line := scanner.Text()

		go func(l string) {
			defer wg.Done()
			processLine(l, stats, alerts)
		}(line)
	}

	// Cerrar canal cuando termine
	go func() {
		wg.Wait()
		close(alerts)
	}()

	// Mostrar alertas
	for alert := range alerts {
		fmt.Println(alert)
		time.Sleep(100 * time.Millisecond)
	}

	stats.Print()
}
