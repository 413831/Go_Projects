# API REST de Gestión de Usuarios

Aplicación de ejemplo en Go para administrar usuarios con funcionalidades avanzadas de seguridad, roles, y gestión de datos personales.

## Características

- ✅ CRUD completo de usuarios (crear, leer, actualizar, borrado lógico)
- ✅ Sistema de roles y permisos
- ✅ Almacenamiento de datos PII (Personally Identifiable Information) con encriptación
- ✅ Control de sesiones con logs
- ✅ Encriptación AES-256-GCM para datos sensibles
- ✅ Hash de contraseñas con bcrypt
- ✅ Logging completo de operaciones
- ✅ Arquitectura escalable con patrones de diseño (Repository, Service)
- ✅ Concurrencia con goroutines, channels y waitgroups
- ✅ Tests unitarios
- ✅ Base de datos PostgreSQL

## Estructura del Proyecto

```
.
├── config/          # Configuración de la aplicación
├── controllers/     # Controladores HTTP
├── database/        # Conexión y migraciones SQL
├── models/          # Modelos de datos
├── repositories/    # Capa de acceso a datos (Repository Pattern)
├── router/          # Configuración de rutas
├── services/        # Lógica de negocio
└── utils/           # Utilidades (encriptación, logging, passwords)
```

## Requisitos

- Go 1.21 o superior
- PostgreSQL (opcional para pruebas, se puede usar mock repository)

## Instalación

1. Clonar el repositorio
2. Instalar dependencias:
```bash
go mod tidy
```

3. Configurar variables de entorno:
   
   **IMPORTANTE**: Copia el archivo `.env.example` a `.env` y configura los valores:
   ```bash
   cp .env.example .env
   ```
   
   Luego edita `.env` con tus valores. **NUNCA** subas el archivo `.env` al repositorio.
   
   Para generar claves seguras:
   ```bash
   # Generar ENCRYPTION_KEY (32 bytes en hexadecimal)
   openssl rand -hex 32
   
   # Generar JWT_SECRET (64 caracteres aleatorios)
   openssl rand -base64 48
   ```
   
   **En producción**, asegúrate de:
   - Definir todas las variables de entorno de seguridad
   - Usar contraseñas fuertes para la base de datos
   - Generar claves de encriptación únicas y seguras
   - Establecer `ENV=production` para habilitar validaciones de seguridad

4. Ejecutar migraciones SQL (ver `database/migrations.sql`)

5. Ejecutar la aplicación:
```bash
go run main.go
```

## Endpoints

### Usuarios

- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users` - Listar usuarios (con paginación: `?page=1&page_size=10`)
- `GET /api/v1/users/{id}` - Obtener usuario por ID
- `PUT /api/v1/users/{id}` - Actualizar usuario
- `DELETE /api/v1/users/{id}` - Borrado lógico de usuario

### Roles

- `POST /api/v1/users/{id}/roles` - Otorgar rol a usuario
- `DELETE /api/v1/users/{id}/roles/{role_id}` - Revocar rol

### Sesiones

- `GET /api/v1/users/{id}/sessions` - Obtener sesiones de un usuario

### Health Check

- `GET /health` - Verificar estado del servicio

## Ejemplos de Uso

### Crear un usuario

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "roles": ["user", "admin"],
    "pii": {
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "1234567890",
      "ssn": "123-45-6789",
      "date_of_birth": "1990-01-01T00:00:00Z"
    }
  }'
```

### Obtener un usuario

```bash
curl http://localhost:8080/api/v1/users/1
```

### Otorgar un rol

```bash
curl -X POST http://localhost:8080/api/v1/users/1/roles \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 1
  }'
```

## Tests

Ejecutar tests unitarios:

```bash
go test ./...
```

Ejecutar tests con cobertura:

```bash
go test -cover ./...
```

## Patrones de Diseño Implementados

- **Repository Pattern**: Separación de la lógica de acceso a datos
- **Service Layer**: Lógica de negocio separada de los controladores
- **Dependency Injection**: Inyección de dependencias para bajo acoplamiento
- **Singleton**: Logger como singleton
- **Factory**: Creación de servicios y repositorios

## Concurrencia

La aplicación utiliza:
- **Goroutines**: Para operaciones paralelas (carga de roles y PII, limpieza de sesiones)
- **Channels**: Para comunicación entre goroutines (limpieza de sesiones)
- **WaitGroups**: Para sincronización de goroutines

## Seguridad

- Contraseñas hasheadas con bcrypt (costo configurable)
- Datos PII encriptados con AES-256-GCM
- Tokens de sesión seguros
- Validación de entrada
- Borrado lógico en lugar de eliminación física

## Seguridad

### Variables de Entorno Requeridas en Producción

En producción (`ENV=production`), las siguientes variables **DEBEN** estar definidas:

- `ENCRYPTION_KEY`: Clave de 32 bytes exactos para encriptación AES-256
- `JWT_SECRET`: Secret para tokens JWT (mínimo 32 caracteres, recomendado 64+)
- `DB_PASSWORD`: Contraseña de la base de datos

La aplicación validará automáticamente que estas variables no sean valores por defecto en producción.

### Generación de Claves Seguras

```bash
# Generar ENCRYPTION_KEY (32 bytes)
openssl rand -hex 32

# Generar JWT_SECRET (64 caracteres)
openssl rand -base64 48

# O usando Python
python3 -c "import secrets; print(secrets.token_hex(32))"  # Para ENCRYPTION_KEY
python3 -c "import secrets; print(secrets.token_urlsafe(48))"  # Para JWT_SECRET
```

### Buenas Prácticas

1. **NUNCA** subas archivos `.env` al repositorio (ya está en `.gitignore`)
2. Usa un gestor de secretos en producción (AWS Secrets Manager, HashiCorp Vault, etc.)
3. Rota las claves periódicamente
4. Usa diferentes claves para cada entorno (desarrollo, staging, producción)
5. Limita el acceso a las variables de entorno solo al personal autorizado

## Notas

- La aplicación está configurada para usar un repositorio mock por defecto (sin base de datos)
- Para usar PostgreSQL, descomentar las líneas en `main.go` y configurar la conexión
- Las claves de encriptación deben ser de 32 bytes exactos
- En desarrollo, se permiten valores por defecto, pero en producción se validan estrictamente

## Licencia

Este es un proyecto de ejemplo para fines educativos.
