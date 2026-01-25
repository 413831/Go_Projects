# Kiosco - API REST con Arquitectura Hexagonal

Este proyecto es una API REST desarrollada en Go que implementa una arquitectura hexagonal (también conocida como arquitectura de puertos y adaptadores) para gestionar items de un kiosco.

## Arquitectura

El proyecto sigue los principios de la arquitectura hexagonal:

- **Dominio** (`internal/domain`): Contiene las entidades y puertos (interfaces) del dominio
- **Aplicación** (`internal/application`): Contiene los casos de uso (servicios de aplicación)
- **Adaptadores de Entrada** (`internal/adapter/input`): Maneja las peticiones HTTP (API REST)
- **Adaptadores de Salida** (`internal/adapter/output`): Se comunica con la API externa que actúa como repositorio

## Estructura del Proyecto

```
kiosco/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── domain/                  # Capa de dominio
│   │   ├── item.go              # Entidad Item
│   │   ├── repository.go        # Puerto (interfaz) del repositorio
│   │   └── errors.go            # Errores del dominio
│   ├── application/             # Capa de aplicación (casos de uso)
│   │   └── item_service.go     # Servicio con los casos de uso CRUD
│   ├── adapter/
│   │   ├── input/               # Adaptadores de entrada
│   │   │   └── http/
│   │   │       ├── handler.go   # Handlers HTTP
│   │   │       └── router.go    # Configuración de rutas
│   │   └── output/              # Adaptadores de salida
│   │       └── api/
│   │           └── item_repository.go  # Cliente HTTP para API externa
│   └── config/
│       └── config.go            # Configuración de la aplicación
├── go.mod
└── README.md
```

## Requisitos

- Go 1.21 o superior

## Instalación

1. Clonar el repositorio
2. Instalar dependencias:
```bash
go mod download
```

## Configuración

Crear un archivo `.env` en la raíz del proyecto (opcional):

```env
SERVER_PORT=8080
EXTERNAL_API_URL=http://localhost:3000/api
```

Si no se proporciona el archivo `.env`, se usarán los valores por defecto:
- `SERVER_PORT`: 8080
- `EXTERNAL_API_URL`: http://localhost:3000/api

## Ejecución

```bash
go run cmd/api/main.go
```

El servidor estará disponible en `http://localhost:8080`

## Endpoints de la API

### Health Check
- **GET** `/health` - Verifica el estado del servidor

### Items

- **GET** `/api/items` - Obtener todos los items
- **GET** `/api/items/{id}` - Obtener un item por ID
- **POST** `/api/items` - Crear un nuevo item
- **PUT** `/api/items/{id}` - Actualizar un item existente
- **DELETE** `/api/items/{id}` - Eliminar un item

## Ejemplos de Uso

### Crear un item

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Coca Cola",
    "description": "Bebida gaseosa",
    "price": 150.50,
    "stock": 100,
    "category": "Bebidas"
  }'
```

### Obtener todos los items

```bash
curl http://localhost:8080/api/items
```

### Obtener un item por ID

```bash
curl http://localhost:8080/api/items/1
```

### Actualizar un item

```bash
curl -X PUT http://localhost:8080/api/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Coca Cola",
    "description": "Bebida gaseosa 500ml",
    "price": 180.00,
    "stock": 80,
    "category": "Bebidas"
  }'
```

### Eliminar un item

```bash
curl -X DELETE http://localhost:8080/api/items/1
```

## Modelo de Datos

### Item

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "price": 0.0,
  "stock": 0,
  "category": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## API Externa

Esta aplicación se comunica con una API externa que actúa como repositorio de datos. La API externa debe implementar los siguientes endpoints:

- `GET /api/items` - Obtener todos los items
- `GET /api/items/{id}` - Obtener un item por ID
- `POST /api/items` - Crear un nuevo item
- `PUT /api/items/{id}` - Actualizar un item
- `DELETE /api/items/{id}` - Eliminar un item

## Ventajas de la Arquitectura Hexagonal

1. **Desacoplamiento**: El dominio no depende de frameworks o tecnologías externas
2. **Testabilidad**: Fácil de testear mediante mocks de los puertos
3. **Flexibilidad**: Fácil cambiar adaptadores (por ejemplo, cambiar de HTTP a gRPC)
4. **Mantenibilidad**: Separación clara de responsabilidades

## Desarrollo

Para agregar nuevas funcionalidades:

1. Agregar entidades en `internal/domain`
2. Definir puertos (interfaces) en `internal/domain`
3. Implementar casos de uso en `internal/application`
4. Crear adaptadores de entrada/salida según sea necesario
