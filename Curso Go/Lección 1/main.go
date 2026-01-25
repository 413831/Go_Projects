package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hola mundo")

	//Números
	entero := 10
	decimal := 3.14
	suma := entero + int(decimal)

	// Texto
	mensaje := "Hola, "
	concatenado := mensaje + "Gentleman"
	enMayuscula := strings.ToUpper(concatenado)

	// Booleanos
	esVerdadero := true

	//Arrays y Slice
	arrayFijo := [3]int{1, 2, 3}
	sliceVariable := []int{4, 5, 6}
	sliceVariable = append(sliceVariable, 7)

	// Mapas
	diccionario := map[string]int{
		"clave1": 1,
		"clave2": 2,
	}

	// Struct
	type Persona struct {
		Nombre string
		Edad   int
	}

	persona := Persona{
		Nombre: "Pepito",
		Edad:   30,
	}

	// Imprimir Resultados
	fmt.Println("Suma", suma)
	fmt.Println("Mensaje", enMayuscula)
	fmt.Println("Array", arrayFijo)
	fmt.Println("Slice", sliceVariable)
	fmt.Println("Mapa", diccionario)
	fmt.Println("Struct", persona)
	fmt.Println("Booleano", esVerdadero)
}

// bool true / false | flag o condicionales
// string cadena de caracteres | para representar texto
// int, int8, int16, int32, int64 | entero | controlas el tamaño de los enteros
// uint, uint8, uint16, uint32, uint64 | entero sin signo | cuando no queremos negativos
// float32, float64 | representar valores nméricos reales | números con punto => 32 | 64 el sistema
// byte === uint8 | trabajar con datos binarios
// rune === int32 | cuando queremos representar un solo caracter que ocupa más de un byte
// complex64, complrex128 | cuando tiene una parte real y una parte imaginaria | N + iN
