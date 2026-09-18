package rs

import (
	"errors"
	"fmt"
	"math/rand"
	"tsp-rs/internal/model"
)

type ResultadoCorrida struct {
	Semilla      int64
	Trayectoria  []int
	Costo        float64
	Estadisticas Estadisticas
	EsFactible   bool
	Historial    []PuntoConvergencia
	Err          error
}

type PuntoConvergencia struct {
	Evaluacion int     `json:"evaluacion"`
	Costo      float64 `json:"costo"`
}

var ErrLoteIncompleto = errors.New("no se pudo completar el lote dentro del límite de intentos")

func CalculaLote(rng *rand.Rand, grafica *model.GraficaTSP, T float64, s []int, N float64, params Parametros, stats *Estadisticas) (float64, []int, error) {
	fS, err := f(grafica, s, N)

	if err != nil {
		return 0, nil, err
	}

	c := 0
	var r float64
	intentos := 0

	for c < params.Lote {
		if intentos >= params.MaxIntentosPorLote {
			return 0, nil, fmt.Errorf("T=%f, %d intentos: %w", T, intentos, ErrLoteIncompleto)
		}
		intentos++

		sPrima := Vecino(rng, s)
		fSPrima, err := f(grafica, sPrima, N)

		if err != nil {
			return 0, nil, err
		}

		if stats != nil {
			stats.EvaluacionesTotales++
		}

		if fSPrima <= fS+T {
			s = sPrima
			fS = fSPrima
			c++
			r += fSPrima

			if stats != nil {
				stats.SolucionesAceptadas++
				if model.EsFactible(grafica, s) {
					stats.SolucionesFactibles++
				} else {
					stats.SolucionesNoFactibles++
				}
			}
		}
	}
	return r / float64(params.Lote), s, nil

}

func PorcentajeAceptados(rng *rand.Rand, grafica *model.GraficaTSP, s []int, T, N float64, iteraciones int) (float64, error) {
	sActual := append([]int(nil), s...)
	fS, err := f(grafica, sActual, N)

	if err != nil {
		return 0, err
	}

	c := 0

	for i := 0; i < iteraciones; i++ {
		sPrima := Vecino(rng, sActual)
		fSPrima, err := f(grafica, sPrima, N)

		if err != nil {
			return 0, err
		}
		if fSPrima <= fS+T {
			c++
			sActual = sPrima
			fS = fSPrima
		}
	}
	return float64(c) / float64(iteraciones), nil

}

func ResolverTSP(semilla int64, grafica *model.GraficaTSP, params Parametros) ([]int, float64, Estadisticas, bool, error) {

	rng := rand.New(rand.NewSource(semilla))

	/*dMax, _, err := model.DistanciaMaxima(*grafica)
	if err != nil {
		return nil, 0, Estadisticas{}, false, fmt.Errorf("calculando distancia máxima: %w", err)
	}*/
	N := model.Normalizador(grafica)

	if N == 0 {
		return nil, 0, Estadisticas{}, false, fmt.Errorf("El normalizador N es 0: porfavor revisa la gráfica de entrada")
	}

	//model.CompletarAristas(grafica, dMax)

	s := SolucionInicialVecinoCercano(rng, grafica, N)

	T0, err := TemperaturaInicial(rng, grafica, s, 8.0, params.Aceptacion, params.EpsilonP, N, params.IteracionesPorcentaje)

	if err != nil {
		return nil, 0, Estadisticas{}, false, fmt.Errorf("Calculando temperatura inicial: %w", err)
	}

	mejorS, mejorCosto, estadisticas, err := AceptacionPorUmbrales(rng, grafica, T0, s, N, params)
	if err != nil {
		return nil, 0, estadisticas, false, fmt.Errorf("Corriendo aceptación por umbrales: %w", err)
	}

	//Barrido final
	mejorS, mejorCosto, err = BusquedaLocal(grafica, mejorS, N)
	if err != nil {
		return nil, 0, estadisticas, false, fmt.Errorf("corriendo barrido final: %w", err)
	}

	mejorEsFactible := model.EsFactible(grafica, mejorS)
	return mejorS, mejorCosto, estadisticas, mejorEsFactible, nil
}

func BusquedaLocal(grafica *model.GraficaTSP, s []int, N float64) ([]int, float64, error) {
	actual := append([]int(nil), s...)
	costoActual, err := f(grafica, actual, N)
	if err != nil {
		return nil, 0, err
	}

	n := len(actual)
	mejorando := true

	for mejorando {
		mejorando = false

		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				candidato := aplicar2opt(actual, i, j)
				costoCandidato, err := f(grafica, candidato, N)
				if err != nil {
					return nil, 0, err
				}

				if costoCandidato < costoActual {
					actual = candidato
					costoActual = costoCandidato
					mejorando = true
				}
			}
		}
	}

	return actual, costoActual, nil
}

func aplicar2opt(s []int, i, j int) []int {
	candidato := make([]int, len(s))
	copy(candidato, s)
	for a, b := i, j; a < b; a, b = a+1, b-1 {
		candidato[a], candidato[b] = candidato[b], candidato[a]
	}
	return candidato
}
