package main

import (
	"fmt"

	"Adapter/adapter"
	"Adapter/geometry"
	"Adapter/renderer"
)

// main demuestra el uso del patrón Adapter para convertir imágenes vectoriales a rasterizadas
//
// PATRÓN ADAPTER:
// El patrón Adapter permite que interfaces incompatibles trabajen juntas. En este caso:
//   - Tenemos una interfaz VectorImage (basada en líneas) que no es compatible
//   - Con una interfaz RasterImage (basada en puntos) que necesitamos usar
//   - El adaptador vectorToRasterAdapter convierte automáticamente las líneas vectoriales
//     en puntos rasterizados, permitiendo que ambas interfaces trabajen juntas
//   - Esto es útil cuando tenemos código legacy o librerías que no podemos modificar
//     pero necesitamos integrar con nuevas funcionalidades
//
// IMPLEMENTACIÓN:
// 1. VectorImage: Define líneas con coordenadas (X1,Y1) -> (X2,Y2)
// 2. RasterImage: Define puntos individuales (X,Y)
// 3. vectorToRasterAdapter: Convierte líneas en puntos usando algoritmos de rasterización
// 4. Caché: Optimiza conversiones repetidas almacenando resultados previamente calculados
//
// ARQUITECTURA MODULAR:
// - geometry: Tipos y funciones geométricas (Point, Line, VectorImage, RasterImage)
// - adapter: Implementación del patrón Adapter (VectorToRasterAdapter)
// - renderer: Funciones de renderizado visual (DrawPoints)
// - utils: Funciones utilitarias (Minmax, Abs)
//
// Crea un rectángulo vectorial, lo convierte a formato raster y procesa los puntos
// La segunda conversión demuestra el uso del caché para optimizar conversiones repetidas
func main() {
	fmt.Println("=== Demostración del Patrón Adapter ===")
	fmt.Println("Conversión de imágenes vectoriales a rasterizadas")
	fmt.Println()

	// Crea un rectángulo vectorial de 6x4
	fmt.Println("1. Creando rectángulo vectorial de 6x4...")
	rc := geometry.NewRectangle(6, 4)
	fmt.Printf("   Líneas vectoriales creadas: %d\n", len(rc.Lines))
	for i, line := range rc.Lines {
		fmt.Printf("   Línea %d: (%d,%d) -> (%d,%d)\n", i+1, line.X1, line.Y1, line.X2, line.Y2)
	}
	fmt.Println()

	// Convierte la imagen vectorial a rasterizada
	fmt.Println("2. Convirtiendo a imagen rasterizada...")
	a := adapter.VectorToRaster(rc)
	fmt.Println()

	// Segunda conversión para demostrar el uso del caché
	fmt.Println("3. Segunda conversión (demostrando caché)...")
	_ = adapter.VectorToRaster(rc)
	fmt.Println()

	// Procesa y muestra los puntos rasterizados
	fmt.Println("4. Generando representación visual:")
	fmt.Print(renderer.DrawPoints(a))

	fmt.Println("=== Fin de la demostración ===")
}
