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

func PesoAumentado(grafica *GraficaTSP, idxU, idxV int, dMax float64) float64 {
	peso := grafica.Matriz[idxU][idxV]
	if !math.IsInf(peso, 1) {
		return peso
	}
	distancia := calcularDistancia(
		grafica.Ciudades[idxU],
		grafica.Ciudades[idxV],
	)

	return distancia * dMax
}

// Revisa si una solución es factible
func EsFactible(grafica *GraficaTSP, path []int) bool {
	for i := 0; i < len(path)-1; i++ {
		idxU := path[i]
		idxV := path[i+1]

		if idxU < 0 || idxU >= len(grafica.MatrizOriginal) ||
			idxV < 0 || idxV >= len(grafica.MatrizOriginal) {
			return false
		}

		if math.IsInf(grafica.MatrizOriginal[idxU][idxV], 1) {
			return false
		}
	}
	return true
}

func EvaluarTrayectoria(grafica *GraficaTSP, path []int, N float64) (float64, error) {
	if len(path) < 2 {
		return 0, fmt.Errorf("La trayectoria necesita al menos 2 ciudaes")
	}
	var suma float64
	for i := 0; i < len(path)-1; i++ {
		idxU := path[i]
		idxV := path[i+1]

		if idxU < 0 || idxU >= len(grafica.Matriz) ||
			idxV < 0 || idxV >= len(grafica.Matriz) {
			return 0, fmt.Errorf(
				"El índice %d o %d está fuera de la matriz",
				idxU,
				idxV,
			)
		}
		suma += PesoAumentado(grafica, idxU, idxV, N)
	}
	return suma, nil
}

func Costo(evaluacion, N float64) float64 {
	if N == 0 {
		return math.Inf(1)
	}
	return evaluacion / N
}
