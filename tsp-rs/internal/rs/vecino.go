package rs

import (
	"math/rand"
)

func Vecino(rng *rand.Rand, s []int) []int {
	n := len(s)
	vecino := make([]int, n)
	copy(vecino, s)

	if n < 3 {
		return vecino
	}
	i := rng.Intn(n)
	j := rng.Intn(n)

	for j == i {
		j = rng.Intn(n)
	}

	if i > j {
		i, j = j, i
	}

	for a, b := i, j; a < b; a, b = a+1, b-1 {
		vecino[a], vecino[b] = vecino[b], vecino[a]
	}
	return vecino
}
