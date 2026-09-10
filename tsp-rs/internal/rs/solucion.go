package rs

import (
	"fmt"
	"math"
	"math/rand"
	"tsp-rs/internal/model"
)

func SolucionInicial(rng *rand.Rand, grafica *GraficaTSP) []int {
	n := len(grafica.Ciudades)
	indices := rng.Perm(n)
	s := make([]int, n+1)
	for i, idx := range indices {
		s[i] = grafica.Ciudades[idx].ID
	}
	return s
}

// f evalúa la función objetivo (costo normalizado) de una solución completa
func f(grafica *GraficaTSP, s []int, N float64) (float64, error) {
	evalucion, err := model.EvaluarTrayectoria(grafica, s, N)
	if err != nil {
		return 0, err
	}
	return model.Costo(evaluacion, N), nil
}
