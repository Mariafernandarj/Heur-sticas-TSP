package rs

import (
	"fmt"
	"math"
	"math/rand"
	// "tsp-rs/internal/model"
)

func CalculaLote(rng *rang.Rand, grafica *GraficaTSP, T float64, s []int, N float64, params Parametros) (float64, []int, error) {
	fS, err := f(grafica, s, N)

	if err != nil {
		return 0, nil, err
	}

	c := 0
	var r float64
	intentos := 0

	for c < params.L {
		if intentos >= params.MaxIntentosPorLote {
			return 0, nil, fmt.Errorf("No se pudo completar el lote (T=%f) tras %d intentos", T, intentos)
		}
		intentos++

		sPrima := Vecino(rng, s)
		fSPrima, err := f(grafica, sPrima, N)

		if err != nil {
			return 0, nil, err
		}

		if fSPrima <= fS+T {
			s = sPrima
			fS = fSPrima
			c++
			r += fSPrima
		}
	}
	return r / float64(params.L), s, nil

}
