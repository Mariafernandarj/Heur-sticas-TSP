package rs

import (
	"errors"
	"math"
	"math/rand"
	"tsp-rs/internal/model"
)

func AceptacionPorUmbrales(rng *rand.Rand, grafica *model.GraficaTSP, T float64, s []int, N float64, params Parametros) ([]int, float64, Estadisticas, error) {

	estadisticas := Estadisticas{}

	mejorS := append([]int(nil), s...)
	mejorCosto, err := f(grafica, s, N)

	if err != nil {
		return nil, 0, estadisticas, err
	}

	p := 0.0
	for T > params.Epsilon {
		q := math.Inf(1)
		lotes := 0
		for p <= q {
			if lotes >= params.MaxLotesPorTemperatura {
				break
			}
			lotes++
			q = p
			pNuevo, sNuevo, err := CalculaLote(rng, grafica, T, s, N, params, &estadisticas)
			if err != nil {
				if errors.Is(err, ErrLoteIncompleto) {
					// El umbral ya es demasiado angosto para seguir
					// aceptando vecinos: la heurística convergió. Nos
					// quedamos con la mejor solución vista hasta ahora en lugar de fallar
					return mejorS, mejorCosto, estadisticas, nil
				}
				return nil, 0, estadisticas, err
			}
			p, s = pNuevo, sNuevo

			costoActual, err := f(grafica, s, N)
			if err != nil {
				return nil, 0, estadisticas, err
			}
			if costoActual < mejorCosto {
				mejorCosto = costoActual
				mejorS = append([]int(nil), s...)
			}
		}
		T *= params.Phi
	}
	return mejorS, mejorCosto, estadisticas, nil
}

func BusquedaBinaria(rng *rand.Rand, grafica *model.GraficaTSP, s []int, T1, T2, Aceptacion, epsilonP, N float64, iteraciones int) (float64, error) {
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
		return BusquedaBinaria(rng, grafica, s, T1, Tm, Aceptacion, epsilonP, N, iteraciones)
	}

	return BusquedaBinaria(rng, grafica, s, Tm, T2, Aceptacion, epsilonP, N, iteraciones)
}

func TemperaturaInicial(rng *rand.Rand, grafica *model.GraficaTSP, s []int, T, Aceptacion, epsilonP, N float64, iteraciones int) (float64, error) {
	p, err := PorcentajeAceptados(rng, grafica, s, T, N, iteraciones)
	if err != nil {
		return 0, err
	}

	if math.Abs(Aceptacion-p) <= epsilonP {
		return T, nil
	}

	var T1, T2 float64
	if p < Aceptacion {
		for p < Aceptacion {
			T *= 2
			p, err = PorcentajeAceptados(rng, grafica, s, T, N, iteraciones)
			if err != nil {
				return 0, err
			}
		}
		T1, T2 = T/2, T
	} else {
		for p > Aceptacion {
			T /= 2
			p, err = PorcentajeAceptados(rng, grafica, s, T, N, iteraciones)
			if err != nil {
				return 0, err
			}
		}
		T1, T2 = T, T*2
	}
	return BusquedaBinaria(rng, grafica, s, T1, T2, Aceptacion, epsilonP, N, iteraciones)
}
