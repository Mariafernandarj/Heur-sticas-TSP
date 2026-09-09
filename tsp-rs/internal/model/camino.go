package model

import (
	"fmt"
	"math"
	"sort"
)

type Arista struct {
	U, V int
	Peso float64
}

const R = 6_373_000.0 // Radio de la Tierra

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

func DistanciaMaxima(grafica GraficaTSP) (float64, Arista, error) {
	mejor := Arista{Peso: math.Inf(-1)}
	for i, fila := range grafica.Matriz {
		for j := i + 1; j < len(fila); j++ {
			peso := fila[j]
			if math.IsInf(peso, 1) {
				continue
			}
			if peso > mejor.Peso {
				mejor = Arista{
					U:    grafica.Ciudades[i].ID,
					V:    grafica.Ciudades[j].ID,
					Peso: peso,
				}
			}
		}
	}
	if math.IsInf(mejor.Peso, -1) {
		return 0, Arista{}, fmt.Errorf("La gráfica no tiene aristas (E está vacío)")
	}
	return mejor.Peso, mejor, nil
}

// El normalizador que suma las k-1 aristas mas pesadas
func Normalizador(grafica *GraficaTSP) float64 {
	n := len(grafica.Matriz)

	pesos := make([]float64, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			peso := grafica.Matriz[i][j]
			if !math.IsInf(peso, 1) {
				pesos = append(pesos, peso)
			}
		}
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(pesos)))

	limite := n - 1
	if limite > len(pesos) {
		limite = len(pesos)
	}

	var N float64
	for _, p := range pesos[:limite] {
		N += p
	}
	return N
}

// función de peso aumentada
func PesoAumentado() {

}
