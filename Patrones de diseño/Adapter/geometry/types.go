package geometry

// Point representa un punto en el espacio 2D
type Point struct {
	X, Y int
}

// Line representa una línea en el espacio 2D con coordenadas de inicio y fin
type Line struct {
	X1, Y1, X2, Y2 int
}

// VectorImage representa una imagen vectorial compuesta por líneas
type VectorImage struct {
	Lines []Line
}

// RasterImage define la interfaz para imágenes rasterizadas
type RasterImage interface {
	GetPoints() []Point
}
