package adapter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"Adapter/geometry"
	"Adapter/utils"
)

// VectorToRasterAdapter implementa el patrón Adapter para convertir imágenes vectoriales a rasterizadas
type VectorToRasterAdapter struct {
	points []geometry.Point
}

// GetPoints implementa el método de la interfaz RasterImage
// Retorna: []Point - slice de puntos que representa la imagen rasterizada
// Este método permite que VectorToRasterAdapter cumpla con la interfaz RasterImage
func (v VectorToRasterAdapter) GetPoints() []geometry.Point {
	return v.points
}

// Caché global para optimizar la conversión de líneas repetidas
var pointCache = map[[16]byte][]geometry.Point{}

// AddLine convierte una línea vectorial en puntos rasterizados y los agrega al adaptador
// Parámetros: line - la línea vectorial a convertir
// La función utiliza un caché para optimizar la conversión de líneas repetidas
// y genera puntos para líneas horizontales, verticales y diagonales
func (v *VectorToRasterAdapter) AddLine(line geometry.Line) {
	// Función hash para generar una clave única basada en el contenido de la línea
	hash := func(obj interface{}) [16]byte {
		bytes, _ := json.Marshal(obj)
		return md5.Sum(bytes)
	}

	h := hash(line)

	// Verifica si la línea ya fue procesada y está en caché
	if pts, ok := pointCache[h]; ok {
		v.points = append(v.points, pts...)
		return
	}

	// Calcula las coordenadas ordenadas de la línea
	left, right := utils.Minmax(line.X1, line.X2)
	top, bottom := utils.Minmax(line.Y1, line.Y2)
	dx := right - left
	dy := bottom - top

	// Genera puntos para líneas verticales (dx = 0)
	if dx == 0 {
		for y := top; y <= bottom; y++ {
			v.points = append(v.points, geometry.Point{X: left, Y: y})
		}
	} else if dy == 0 { // Genera puntos para líneas horizontales (dy = 0)
		for x := left; x <= right; x++ {
			v.points = append(v.points, geometry.Point{X: x, Y: top})
		}
	} else { // Genera puntos para líneas diagonales usando el algoritmo de Bresenham
		// Algoritmo de Bresenham para líneas diagonales
		x0, y0 := line.X1, line.Y1
		x1, y1 := line.X2, line.Y2

		dx := utils.Abs(x1 - x0)
		dy := utils.Abs(y1 - y0)

		var sx, sy int
		if x0 < x1 {
			sx = 1
		} else {
			sx = -1
		}
		if y0 < y1 {
			sy = 1
		} else {
			sy = -1
		}

		err := dx - dy
		x, y := x0, y0

		for {
			v.points = append(v.points, geometry.Point{X: x, Y: y})

			if x == x1 && y == y1 {
				break
			}

			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x += sx
			}
			if e2 < dx {
				err += dx
				y += sy
			}
		}
	}

	// Guarda los puntos generados en el caché para futuras consultas
	pointCache[h] = v.points
	fmt.Println("generated", len(v.points), "points")
}

// VectorToRaster convierte una imagen vectorial en una imagen rasterizada usando el patrón Adapter
// Parámetros: vi - puntero a la imagen vectorial a convertir
// Retorna: RasterImage - interfaz que representa la imagen rasterizada
// Esta función es el punto de entrada principal del adaptador que procesa todas las líneas
// de la imagen vectorial y las convierte en puntos rasterizados
func VectorToRaster(vi *geometry.VectorImage) geometry.RasterImage {
	adapter := VectorToRasterAdapter{}

	// Procesa cada línea de la imagen vectorial
	for _, line := range vi.Lines {
		adapter.AddLine(line)
	}

	return adapter // como RasterImage
}

// ClearCache limpia el caché de puntos
func ClearCache() {
	pointCache = make(map[[16]byte][]geometry.Point)
}

// HasCacheEntries verifica si el caché tiene entradas
func HasCacheEntries() bool {
	return len(pointCache) > 0
}

// CacheSize devuelve el número de entradas en el caché
func CacheSize() int {
	return len(pointCache)
}
