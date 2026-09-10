package rs

import (
	// "fmt"
	"math"
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

func BusquedaBinaria(rng *rand.Rand, grafica *GraficaTSP, s []int, T1, T2, Aceptacion, epsilonP, N float64, iteraciones int) (float64, error) {
	Tm := (T1 + T2) / 2
	if T2-T1 < epsilonP {
		return Tm, nil
	}

	p, err := PorcentajeAceptados(rng, grafica, s, Tm, N, iteraciones)

	if err != nil {
		return 0, err
	}

	if math.Abs(Aceptacion-p) < epsilonP {
		return Tm, nil
	}

	if p > Aceptacion {
		return BusquedaBInaria(rng, grafica, s, T1, Tm, Aceptacion, epsilonP, N, iteraciones)
	}

	return BusquedaBinaria(rng, grafica, s, Tm, T2, Aceptacion, epsilonP, N, iteraciones)
}
