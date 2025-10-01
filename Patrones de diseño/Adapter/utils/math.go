package utils

// Minmax devuelve los valores mínimo y máximo de dos enteros en orden ascendente
// Parámetros: a, b - los dos enteros a comparar
// Retorna: (min, max) - el valor menor y el valor mayor respectivamente
func Minmax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

// Abs devuelve el valor absoluto de un entero
// Parámetros: x - el entero del cual obtener el valor absoluto
// Retorna: int - el valor absoluto de x
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
