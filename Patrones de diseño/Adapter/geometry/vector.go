package geometry

// NewRectangle crea una nueva imagen vectorial rectangular con las dimensiones especificadas
// Parámetros: width, height - ancho y alto del rectángulo
// Retorna: *VectorImage - puntero a una imagen vectorial que representa un rectángulo
// La función ajusta las coordenadas restando 1 para crear líneas que formen un rectángulo cerrado
func NewRectangle(width, height int) *VectorImage {
	width -= 1
	height -= 1

	return &VectorImage{[]Line{
		{0, 0, width, 0},           // línea superior
		{0, 0, 0, height},          // línea izquierda
		{width, 0, width, height},  // línea derecha
		{0, height, width, height}, // línea inferior
	}}
}
