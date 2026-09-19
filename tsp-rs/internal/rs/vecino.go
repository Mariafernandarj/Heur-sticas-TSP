package rs

import (
	"math/rand"
)

/* Genera una nueva solución candidata (vecino) aplicando un movimiento 2-opt (inversión)
 *sobre la secuencia actual 's', utilizando el generador aleatorio 'rng'
 *para garantizar Hilo-Seguridad (Thread-Safety)
 */
func Vecino(rng *rand.Rand, s []int) []int {
	n := len(s)
	vecino := make([]int, n)
	copy(vecino, s) // Duplica la solución actual para no mutar el slice original

	// Se requieren al menos 3 elementos para realizar una inversión con sentido en un tour
	if n < 3 {
		return vecino
	}

	// Selecciona dos índices aleatorios distintos
	i := rng.Intn(n)
	j := rng.Intn(n)

	for j == i {
		j = rng.Intn(n)
	}

	if i > j {
		i, j = j, i
	}

	// Invierte el subsegmento de la solución entre los índices [i, j]
	for a, b := i, j; a < b; a, b = a+1, b-1 {
		vecino[a], vecino[b] = vecino[b], vecino[a]
	}
	return vecino
}
