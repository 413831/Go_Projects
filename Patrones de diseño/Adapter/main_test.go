package main

import (
	"testing"

	"Adapter/adapter"
	"Adapter/geometry"
	"Adapter/renderer"
	"Adapter/utils"
)

// TestMinmax verifica que la función Minmax devuelva correctamente los valores mínimo y máximo
func TestMinmax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected [2]int
	}{
		{"a menor que b", 3, 7, [2]int{3, 7}},
		{"a mayor que b", 9, 4, [2]int{4, 9}},
		{"a igual a b", 5, 5, [2]int{5, 5}},
		{"números negativos", -3, -1, [2]int{-3, -1}},
		{"uno negativo uno positivo", -2, 3, [2]int{-2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := utils.Minmax(tt.a, tt.b)
			if min != tt.expected[0] || max != tt.expected[1] {
				t.Errorf("Minmax(%d, %d) = (%d, %d), esperado (%d, %d)",
					tt.a, tt.b, min, max, tt.expected[0], tt.expected[1])
			}
		})
	}
}

// TestAbs verifica que la función Abs devuelva correctamente el valor absoluto
func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"número positivo", 5, 5},
		{"número negativo", -3, 3},
		{"cero", 0, 0},
		{"número grande positivo", 1000, 1000},
		{"número grande negativo", -1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs(%d) = %d, esperado %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNewRectangle verifica que la función NewRectangle cree correctamente un rectángulo vectorial
func TestNewRectangle(t *testing.T) {
	tests := []struct {
		name           string
		width, height  int
		expectedLines  int
		expectedCoords [][4]int // [X1, Y1, X2, Y2]
	}{
		{
			name:          "rectángulo 3x2",
			width:         3,
			height:        2,
			expectedLines: 4,
			expectedCoords: [][4]int{
				{0, 0, 2, 0}, // línea superior
				{0, 0, 0, 1}, // línea izquierda
				{2, 0, 2, 1}, // línea derecha
				{0, 1, 2, 1}, // línea inferior
			},
		},
		{
			name:          "rectángulo 1x1",
			width:         1,
			height:        1,
			expectedLines: 4,
			expectedCoords: [][4]int{
				{0, 0, 0, 0}, // línea superior
				{0, 0, 0, 0}, // línea izquierda
				{0, 0, 0, 0}, // línea derecha
				{0, 0, 0, 0}, // línea inferior
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect := geometry.NewRectangle(tt.width, tt.height)

			if len(rect.Lines) != tt.expectedLines {
				t.Errorf("NewRectangle(%d, %d) creó %d líneas, esperado %d",
					tt.width, tt.height, len(rect.Lines), tt.expectedLines)
			}

			for i, line := range rect.Lines {
				expected := tt.expectedCoords[i]
				if line.X1 != expected[0] || line.Y1 != expected[1] ||
					line.X2 != expected[2] || line.Y2 != expected[3] {
					t.Errorf("Línea %d: (%d,%d)->(%d,%d), esperado (%d,%d)->(%d,%d)",
						i, line.X1, line.Y1, line.X2, line.Y2,
						expected[0], expected[1], expected[2], expected[3])
				}
			}
		})
	}
}

// TestVectorToRasterAdapter verifica que el adaptador convierta correctamente líneas en puntos
func TestVectorToRasterAdapter(t *testing.T) {
	// Crear un rectángulo simple de 2x2
	rect := geometry.NewRectangle(2, 2)
	adapter := adapter.VectorToRaster(rect)
	points := adapter.GetPoints()

	// Verificar que se generaron puntos
	if len(points) == 0 {
		t.Error("El adaptador no generó ningún punto")
	}

	// Verificar que los puntos están dentro del rango esperado
	for _, point := range points {
		if point.X < 0 || point.X >= 2 || point.Y < 0 || point.Y >= 2 {
			t.Errorf("Punto fuera de rango: (%d, %d)", point.X, point.Y)
		}
	}

	// Verificar que hay al menos algunos puntos esperados para un rectángulo 2x2
	expectedMinPoints := 4 // Mínimo para un rectángulo 2x2
	if len(points) < expectedMinPoints {
		t.Errorf("Se generaron %d puntos, esperado al menos %d", len(points), expectedMinPoints)
	}
}

// TestDrawPoints verifica que la función DrawPoints genere una salida válida
func TestDrawPoints(t *testing.T) {
	// Crear un rectángulo pequeño para testing
	rect := geometry.NewRectangle(2, 2)
	adapter := adapter.VectorToRaster(rect)
	result := renderer.DrawPoints(adapter)

	// Verificar que la salida no esté vacía
	if result == "" {
		t.Error("DrawPoints devolvió una cadena vacía")
	}

	// Verificar que contiene información esperada
	expectedStrings := []string{
		"Imagen rasterizada",
		"Total de puntos:",
		"*", // Debe contener al menos un punto activo
	}

	for _, expected := range expectedStrings {
		if !contains(result, expected) {
			t.Errorf("DrawPoints no contiene '%s' en la salida: %s", expected, result)
		}
	}
}

// TestAddLineHorizontal verifica que AddLine procese correctamente líneas horizontales
func TestAddLineHorizontal(t *testing.T) {
	adapter := &adapter.VectorToRasterAdapter{}
	line := geometry.Line{X1: 0, Y1: 0, X2: 3, Y2: 0} // Línea horizontal
	adapter.AddLine(line)

	points := adapter.GetPoints()
	if len(points) == 0 {
		t.Error("addLine no generó puntos para línea horizontal")
	}

	// Verificar que todos los puntos están en la misma fila (Y=0)
	for _, point := range points {
		if point.Y != 0 {
			t.Errorf("Punto de línea horizontal tiene Y=%d, esperado 0", point.Y)
		}
	}

	// Verificar que hay al menos 4 puntos (0,0), (1,0), (2,0), (3,0)
	if len(points) < 4 {
		t.Errorf("Línea horizontal generó %d puntos, esperado al menos 4", len(points))
	}
}

// TestAddLineVertical verifica que AddLine procese correctamente líneas verticales
func TestAddLineVertical(t *testing.T) {
	adapter := &adapter.VectorToRasterAdapter{}
	line := geometry.Line{X1: 0, Y1: 0, X2: 0, Y2: 3} // Línea vertical
	adapter.AddLine(line)

	points := adapter.GetPoints()
	if len(points) == 0 {
		t.Error("addLine no generó puntos para línea vertical")
	}

	// Verificar que todos los puntos están en la misma columna (X=0)
	for _, point := range points {
		if point.X != 0 {
			t.Errorf("Punto de línea vertical tiene X=%d, esperado 0", point.X)
		}
	}

	// Verificar que hay al menos 4 puntos (0,0), (0,1), (0,2), (0,3)
	if len(points) < 4 {
		t.Errorf("Línea vertical generó %d puntos, esperado al menos 4", len(points))
	}
}

// TestCacheFunctionality verifica que el caché funcione correctamente
func TestCacheFunctionality(t *testing.T) {
	// Limpiar el caché antes del test
	adapter.ClearCache()

	rect := geometry.NewRectangle(2, 2)

	// Primera conversión
	adapter1 := adapter.VectorToRaster(rect)
	points1 := adapter1.GetPoints()

	// Verificar que el caché no esté vacío después de la primera conversión
	if !adapter.HasCacheEntries() {
		t.Error("El caché debería contener entradas después de la primera conversión")
	}

	// Verificar que se generaron puntos
	if len(points1) == 0 {
		t.Error("La primera conversión no generó puntos")
	}

	// Verificar que el caché contiene las líneas esperadas (4 líneas para un rectángulo)
	if adapter.CacheSize() < 4 {
		t.Errorf("El caché debería contener al menos 4 entradas (una por línea), tiene %d", adapter.CacheSize())
	}
}

// TestDrawPointsEmpty verifica el comportamiento con una imagen vacía
func TestDrawPointsEmpty(t *testing.T) {
	// Crear un adaptador vacío
	adapter := &adapter.VectorToRasterAdapter{}
	result := renderer.DrawPoints(adapter)

	expected := "No hay puntos para dibujar\n"
	if result != expected {
		t.Errorf("DrawPoints con imagen vacía devolvió: '%s', esperado: '%s'", result, expected)
	}
}

// Función auxiliar para verificar si una cadena contiene una subcadena
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr))))
}
