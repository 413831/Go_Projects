package renderer

import (
	"fmt"

	"Adapter/geometry"
)

// DrawPoints procesa una imagen raster y genera una representación visual de los puntos
// Parámetros: owner - implementación de RasterImage que contiene los puntos a procesar
// Retorna: string - representación visual de la imagen rasterizada
// La función encuentra las coordenadas X e Y máximas y crea una matriz de caracteres
func DrawPoints(owner geometry.RasterImage) string {
	maxX, maxY := 0, 0
	points := owner.GetPoints()

	if len(points) == 0 {
		return "No hay puntos para dibujar\n"
	}

	// Encuentra las coordenadas máximas en X e Y
	for _, pixel := range points {
		if pixel.X > maxX {
			maxX = pixel.X
		}
		if pixel.Y > maxY {
			maxY = pixel.Y
		}
	}

	maxX += 1
	maxY += 1

	// Crea un mapa para verificar rápidamente si un punto existe
	pointMap := make(map[geometry.Point]bool)
	for _, point := range points {
		pointMap[point] = true
	}

	// Genera la representación visual
	var result string
	result += fmt.Sprintf("Imagen rasterizada (%dx%d):\n", maxX, maxY)

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if pointMap[geometry.Point{X: x, Y: y}] {
				result += "*" // Punto activo
			} else {
				result += "." // Espacio vacío
			}
		}
		result += "\n"
	}

	result += fmt.Sprintf("Total de puntos: %d\n", len(points))
	return result
}
