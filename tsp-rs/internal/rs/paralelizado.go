package rs

import (
	"sync"
	"tsp-rs/internal/model"
)

/* Ejecuta de manera concurrente la heurística del TSP para una lista de semillas
 *utilizando goroutines y un sync.WaitGroup
 * cada ejecución corre de forma aislada en su propio hilo ligero
 */
func ResolverTspParalelo(grafica *model.GraficaTSP, params Parametros, semillas []int64) []ResultadoCorrida {

	resultados := make([]ResultadoCorrida, len(semillas))

	var wg sync.WaitGroup
	for i, semilla := range semillas {
		wg.Add(1)
		// Lanza una goroutine por cada semilla, pasando 'i' y 'semilla' por valor
		// para garantizar un ámbito (scope) independiente en cada iteración
		go func(i int, semilla int64) {
			defer wg.Done()

			// Realiza una copia profunda de la gráfica para garantizar
			// la independencia de memoria entre goroutines y evitar Data Races
			graficaLocal := model.CopiarGrafica(*grafica)
			trayectoria, costo, estads, esFactible, err := ResolverTSP(semilla, graficaLocal, params)

			// Guarda la estructura del resultado en la posición preasignada correspondiente
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

/* Filtra los resultados paralelos descartando aquellos con errores
 * y devuelve la corrida con el menor costo total alcanzado
 */
func MejorCorrida(resultados []ResultadoCorrida) (ResultadoCorrida, bool) {
	var mejor ResultadoCorrida
	encontrada := false

	// Itera sobre los resultados para seleccionar el óptimo global de la ejecución concurrente
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
