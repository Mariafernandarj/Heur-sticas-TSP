package model

import (
	//"fmt"
	"math"
)

const R = 6_373_000.0 // Radio de la Tierra

// Representación de una fila en la tabla 'cities'
type Ciudad struct {
	ID        int
	Nombre    string
	Pais      string
	Poblacion int
	Latitud   float64
	Longitud  float64
}

// Recibe latitud/longitud en grados y lo devuelve la distancia en metros
func DistanciaNatural(latU, longU, latV, longV float64) float64 {
	latUR := latU * math.Pi / 180
	longUR := longU * math.Pi / 180
	latVR := latV * math.Pi / 180
	longVR := longV * math.Pi / 180

	distanciaLat := latVR - latUR
	distanciaLong := longVR - longUR

	A := math.Pow(math.Sin(distanciaLat/2), 2) +
		math.Cos(latUR)*math.Cos(latVR)*math.Pow(math.Sin(distanciaLong/2), 2)

	C := 2 * math.Atan2(math.Sqrt(A), math.Sqrt(1-A))

	return R * C
}
