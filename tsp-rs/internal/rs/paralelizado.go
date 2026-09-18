package rs

import (
	"sync"
	"tsp-rs/internal/model"
)

func ResolverTspParalelo(grafica *model.GraficaTSP, params Parametros, semillas []int64) []ResultadoCorrida {

	resultados := make([]ResultadoCorrida, len(semillas))

	var wg sync.WaitGroup
	for i, semilla := range semillas {
		wg.Add(1)
		go func(i int, semilla int64) {
			defer wg.Done()

			graficaLocal := model.CopiarGrafica(*grafica)
			trayectoria, costo, estads, esFactible, err := ResolverTSP(semilla, graficaLocal, params)

			resultados[i] = ResultadoCorrida{
				Semilla:      semilla,
				Trayectoria:  trayectoria,
				Costo:        costo,
				Estadisticas: estads,
				EsFactible:   esFactible,
				Historial:    estads.Historial,
				Err:          err,
			}
		}(i, semilla)
	}
	wg.Wait()
	return resultados
}

func MejorCorrida(resultados []ResultadoCorrida) (ResultadoCorrida, bool) {
	var mejor ResultadoCorrida
	encontrada := false

	for _, r := range resultados {
		if r.Err != nil {
			continue
		}
		if !encontrada || r.Costo < mejor.Costo {
			mejor = r
			encontrada = true
		}
	}
	return mejor, encontrada
}
