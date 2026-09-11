package rs

import (
	"math/rand"
	"tsp-rs/internal/model"
)

func SolucionInicial(rng *rand.Rand, grafica *model.GraficaTSP) []int {
	n := len(grafica.Ciudades)
	indices := rng.Perm(n)
	s := make([]int, n)
	for i, idx := range indices {
		s[i] = idx
	}
	return s
}

// f evalúa la función objetivo (costo normalizado) de una solución completa
func f(grafica *model.GraficaTSP, s []int, N float64) (float64, error) {
	evaluacion, err := model.EvaluarTrayectoria(grafica, s, N)
	if err != nil {
		return 0, err
	}
	return model.Costo(evaluacion, N), nil
}
