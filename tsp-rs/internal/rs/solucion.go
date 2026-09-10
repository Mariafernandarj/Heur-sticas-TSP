package rs

import (
	"fmt"
	"math"
	"math/rand"
)

func SolucionInicial(rng *rand.Rand, grafica *GraficaTSP) []int {
	n := len(grafica.Ciudades)
	indices := rng.Perm(n)
	s := make([]int, n+1)
	for i, idx := range indices {
		s[i] = grafica.Ciudades[idx].ID
	}
	s[n] = s[0]
	return s
}
