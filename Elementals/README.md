# Elementals - Proyectos de Aprendizaje de Go

Esta carpeta contiene una colección de proyectos de ejemplo en Go que demuestran diferentes conceptos fundamentales del lenguaje. Cada proyecto está diseñado para enseñar conceptos específicos de programación en Go.

## 📁 Estructura de Proyectos

### 1. Hello_World
**Concepto:** Estructuras básicas y slices
- **Archivo principal:** `main.go`
- **Descripción:** Demuestra la creación de estructuras simples y el uso de slices vacíos
- **Conceptos clave:** Structs, slices, len()

### 2. api_rest_basic
**Concepto:** Servidor HTTP REST básico
- **Archivos:** `main.go`, `server/server.go`, `server/handlers.go`, `server/routes.go`
- **Descripción:** Implementa un servidor HTTP simple con endpoints para manejar países
- **Conceptos clave:** HTTP server, JSON encoding/decoding, routing, graceful shutdown
- **Endpoints:**
  - `GET /` - Página de bienvenida
  - `GET /countries` - Obtener lista de países
  - `POST /countries` - Agregar nuevo país

### 3. cards
**Concepto:** Tipos personalizados y métodos
- **Archivos:** `main.go`, `deck.go`, `deck_test.go`
- **Descripción:** Sistema de cartas con operaciones como mezclar, repartir y guardar/cargar
- **Conceptos clave:** Custom types, methods, file I/O, testing
- **Funcionalidades:**
  - Crear mazo de cartas
  - Mezclar cartas
  - Repartir cartas
  - Guardar/cargar desde archivo
  - Tests unitarios

### 4. channels
**Concepto:** Goroutines y canales
- **Archivo principal:** `main.go`
- **Descripción:** Verificación concurrente del estado de sitios web usando goroutines
- **Conceptos clave:** Goroutines, channels, concurrency, HTTP requests
- **Funcionalidad:** Monitorea múltiples URLs de forma concurrente

### 5. exercise_slices
**Concepto:** Slices y bucles
- **Archivos:** `main.go`, `slice_test.go`
- **Descripción:** Ejercicio básico con slices para determinar números pares e impares
- **Conceptos clave:** Slices, range loops, modulo operator
- **Tests:** Incluye tests para verificar la lógica de detección de pares/impares

### 6. interfaces
**Concepto:** Interfaces básicas
- **Archivos:** `main.go`, `bot_test.go`
- **Descripción:** Demuestra el uso de interfaces con diferentes tipos de bots
- **Conceptos clave:** Interfaces, method implementation, polymorphism
- **Tests:** Verifica que los bots implementen correctamente la interfaz

### 7. interfaces_http
**Concepto:** Interfaces con HTTP
- **Archivo principal:** `main.go`
- **Descripción:** Implementación de la interfaz io.Writer para logging personalizado
- **Conceptos clave:** io.Writer interface, HTTP responses, custom writers

### 8. interfaces2
**Concepto:** Interfaces con formas geométricas
- **Archivos:** `main.go`, `shape_test.go`
- **Descripción:** Cálculo de áreas usando interfaces para diferentes formas
- **Conceptos clave:** Interfaces, geometric calculations, polymorphism
- **Tests:** Verifica cálculos de área para triángulos y cuadrados

### 9. interfaces3
**Concepto:** Interfaces para manipulación de archivos
- **Archivo principal:** `main.go`
- **Archivo de datos:** `myfile.txt`
- **Descripción:** Uso de interfaces io.Reader e io.Writer para copiar archivos
- **Conceptos clave:** io.Reader, io.Writer, file operations, command line arguments

### 10. map
**Concepto:** Mapas (dictionaries)
- **Archivos:** `main.go`, `color_test.go`
- **Descripción:** Uso de mapas para almacenar códigos de colores hexadecimales
- **Conceptos clave:** Maps, key-value pairs, iteration
- **Tests:** Verifica operaciones básicas de mapas

### 11. structs
**Concepto:** Estructuras y punteros
- **Archivos:** `main.go`, `person_test.go`
- **Descripción:** Estructuras anidadas y métodos con punteros
- **Conceptos clave:** Structs, pointers, methods, nested structs
- **Tests:** Verifica creación y modificación de estructuras

## 🚀 Cómo Ejecutar los Proyectos

Cada proyecto es independiente y puede ejecutarse por separado:

```bash
# Navegar a cualquier proyecto
cd nombre_del_proyecto

# Ejecutar el programa
go run main.go

# Ejecutar tests (si están disponibles)
go test

# Ejecutar tests con verbose
go test -v
```

## 📋 Requisitos

- Go 1.19 o superior
- Conexión a internet (para algunos proyectos que hacen peticiones HTTP)

## 🧪 Testing

Los siguientes proyectos incluyen tests unitarios:
- `cards` - Tests para el sistema de cartas
- `exercise_slices` - Tests para detección de pares/impares
- `interfaces` - Tests para bots
- `interfaces2` - Tests para formas geométricas
- `map` - Tests para operaciones de mapas
- `structs` - Tests para estructuras

## 📚 Conceptos de Go Cubiertos

1. **Tipos básicos y variables**
2. **Estructuras (structs)**
3. **Slices y arrays**
4. **Mapas (maps)**
5. **Funciones y métodos**
6. **Punteros**
7. **Interfaces**
8. **Goroutines y canales**
9. **Manejo de archivos**
10. **HTTP y JSON**
11. **Testing**
12. **Manejo de errores**

## 🎯 Objetivo de Aprendizaje

Esta colección está diseñada para proporcionar una base sólida en Go, cubriendo desde conceptos básicos hasta características avanzadas como concurrencia e interfaces. Cada proyecto puede ser estudiado independientemente o como parte de una progresión de aprendizaje.

## 📝 Notas

- Todos los archivos están completamente documentados en español
- Los tests están diseñados para ser educativos y demostrar buenas prácticas
- Los proyectos están organizados por complejidad creciente
- Se incluyen ejemplos de código comentado para facilitar el aprendizaje
