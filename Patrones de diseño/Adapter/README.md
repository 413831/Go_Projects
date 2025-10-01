# Patrón Adapter - Conversión de Imágenes Vectoriales a Rasterizadas

Este proyecto implementa el patrón de diseño Adapter para convertir imágenes vectoriales (basadas en líneas) a imágenes rasterizadas (basadas en puntos).

## Arquitectura Modular

El proyecto está organizado en módulos separados para mejorar la mantenibilidad y organización del código:

### 📁 Estructura del Proyecto

```
Adapter/
├── main.go                 # Punto de entrada principal
├── main_test.go           # Tests del sistema completo
├── go.mod                 # Módulo Go
├── README.md              # Documentación
├── utils/                 # Funciones utilitarias
│   └── math.go           # Funciones matemáticas (Minmax, Abs)
├── geometry/              # Tipos y funciones geométricas
│   ├── types.go          # Definiciones de tipos (Point, Line, VectorImage, RasterImage)
│   └── vector.go         # Funciones para imágenes vectoriales (NewRectangle)
├── adapter/               # Implementación del patrón Adapter
│   └── vector_to_raster.go # Adaptador VectorToRasterAdapter
└── renderer/              # Funciones de renderizado
    └── draw.go           # Función DrawPoints para visualización
```

### 🔧 Módulos

#### **utils** - Funciones Utilitarias
- `Minmax(a, b int) (int, int)`: Devuelve valores mínimo y máximo ordenados
- `Abs(x int) int`: Calcula el valor absoluto

#### **geometry** - Geometría
- **Tipos**:
  - `Point`: Punto 2D con coordenadas X, Y
  - `Line`: Línea 2D con coordenadas de inicio y fin
  - `VectorImage`: Imagen vectorial compuesta por líneas
  - `RasterImage`: Interfaz para imágenes rasterizadas
- **Funciones**:
  - `NewRectangle(width, height int) *VectorImage`: Crea rectángulo vectorial

#### **adapter** - Patrón Adapter
- `VectorToRasterAdapter`: Adaptador que convierte líneas en puntos
- `VectorToRaster(vi *VectorImage) RasterImage`: Función principal de conversión
- `AddLine(line Line)`: Convierte una línea en puntos rasterizados
- Funciones de caché: `ClearCache()`, `HasCacheEntries()`, `CacheSize()`

#### **renderer** - Renderizado
- `DrawPoints(owner RasterImage) string`: Genera representación visual ASCII

## 🚀 Uso

### Ejecutar el programa:
```bash
go run main.go
```

### Ejecutar tests:
```bash
go test -v
```

### Ejecutar tests con cobertura:
```bash
go test -v -cover
```

## 🎯 Patrón Adapter

### ¿Qué es?
El patrón Adapter permite que interfaces incompatibles trabajen juntas. En este caso:
- **VectorImage** (basada en líneas) ↔ **RasterImage** (basada en puntos)
- El adaptador convierte automáticamente las líneas vectoriales en puntos rasterizados

### ¿Cuándo usarlo?
- Integrar código legacy o librerías sin modificar
- Hacer compatibles interfaces que no fueron diseñadas para trabajar juntas
- Crear una capa de abstracción entre sistemas diferentes

### Implementación en este proyecto:
1. **VectorImage**: Define líneas con coordenadas (X1,Y1) → (X2,Y2)
2. **RasterImage**: Define puntos individuales (X,Y)
3. **VectorToRasterAdapter**: Convierte líneas en puntos usando algoritmos de rasterización
4. **Caché**: Optimiza conversiones repetidas almacenando resultados

## 🧪 Testing

El proyecto incluye tests exhaustivos que cubren:
- ✅ Funciones utilitarias (Minmax, Abs)
- ✅ Creación de rectángulos vectoriales
- ✅ Conversión vectorial a rasterizada
- ✅ Renderizado visual
- ✅ Líneas horizontales, verticales y diagonales
- ✅ Funcionalidad del caché
- ✅ Casos edge (imágenes vacías)

## 🔍 Algoritmos Implementados

### Rasterización de Líneas
- **Líneas horizontales**: Generación directa de puntos en Y constante
- **Líneas verticales**: Generación directa de puntos en X constante  
- **Líneas diagonales**: Algoritmo de Bresenham para líneas suaves

### Optimización
- **Caché de puntos**: Evita recálculos de líneas repetidas usando hash MD5
- **Conversión eficiente**: Procesa solo líneas nuevas

## 📊 Ejemplo de Salida

```
=== Demostración del Patrón Adapter ===
Conversión de imágenes vectoriales a rasterizadas

1. Creando rectángulo vectorial de 6x4...
   Líneas vectoriales creadas: 4
   Línea 1: (0,0) -> (5,0)
   Línea 2: (0,0) -> (0,3)
   Línea 3: (5,0) -> (5,3)
   Línea 4: (0,3) -> (5,3)

2. Convirtiendo a imagen rasterizada...
generated 6 points
generated 10 points
generated 14 points
generated 20 points

3. Segunda conversión (demostrando caché)...

4. Generando representación visual:
Imagen rasterizada (6x4):
******
*....*
*....*
******
Total de puntos: 20
=== Fin de la demostración ===
```

## 🛠️ Requisitos

- Go 1.19 o superior
- Módulos Go habilitados

## 📝 Notas de Desarrollo

- El código está completamente documentado en español
- Sigue las mejores prácticas de Go
- Incluye manejo de errores y casos edge
- Optimizado para rendimiento con caché
- Arquitectura modular y extensible
