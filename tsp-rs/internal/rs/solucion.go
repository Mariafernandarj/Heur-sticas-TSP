package rs

import (
	"math"
	"math/rand"
	"tsp-rs/internal/model"
)

type Estadisticas struct {
	SolucionesAceptadas   int
	SolucionesFactibles   int
	SolucionesNoFactibles int
	TopesPorMaxLotes      int
	EvaluacionesTotales   int
	Historial             []PuntoConvergencia
}

func SolucionInicial(rng *rand.Rand, grafica *model.GraficaTSP) []int {
	n := len(grafica.Ciudades)
	indices := rng.Perm(n)
	s := make([]int, n)
	for i, idx := range indices {
		s[i] = idx
	}
	return s
}

func SolucionInicialVecinoCercano(rng *rand.Rand, grafica *model.GraficaTSP, N float64) []int {
	n := len(grafica.Ciudades)
	visitado := make([]bool, n)

	actual := rng.Intn(n)
	visitado[actual] = true

	path := make([]int, 0, n)
	path = append(path, grafica.Ciudades[actual].ID)

	for len(path) < n {
		mejorIdx := -1
		mejorPeso := math.Inf(1)

		for j := 0; j < n; j++ {
			if visitado[j] {
				continue
			}
			peso := model.PesoAumentado(grafica, actual, j, N)
			if peso < mejorPeso {
				mejorPeso = peso
				mejorIdx = j
			}
		}

		visitado[mejorIdx] = true
		path = append(path, grafica.Ciudades[mejorIdx].ID)
		actual = mejorIdx
	}

	return path
}

// f evalúa la función objetivo (costo normalizado) de una solución completa
func f(grafica *model.GraficaTSP, s []int, N float64) (float64, error) {
	evaluacion, err := model.EvaluarTrayectoria(grafica, s, N)
	if err != nil {
		return 0, err
	}
	return model.Costo(evaluacion, N), nil
}

func (e Estadisticas) PorcentajeFactibles() float64 {
	if e.SolucionesAceptadas == 0 {
		return 0
	}
	return float64(e.SolucionesFactibles) / float64(e.SolucionesAceptadas)
}
