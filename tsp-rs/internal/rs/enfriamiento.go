package rs

import (
// "fmt"
// "math"
// "math/rand"
// "tsp-rs/internal/model"
)

func AceptacionPorUmbrales(rng *rand.Rand, grafica *GraficaTSP, T float64, s []int, N float64, params Parametros) ([]int, float64, error) {

	mejorS := append([]int(nil), s...)
	mejorCosto, err := f(grafica, s, N)

	if err != nil {
		return nil, 0, err
	}

	p := 0.0
	for T > params.Epsilon {
		q := math.Inf(1)
		for p <= q {
			q = p
			var err error
			p, s, err = CalculaLote(rng, grafica, T, s, N, params)
			if err != nil {
				return nil, 0, err
			}

			costoActual, err := f(grafica, s, N)
			if err != nil {
				return nil, 0, err
			}
			if costoActual < mejorCosto {
				mejorCosto = costoActual
				mejorS = append([]int(nil), s...)
			}
		}
		T *= params.Phi
	}
	return mejorS, mejorCosto, nil
}
