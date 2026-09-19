package rs

import (
	"errors"
	"fmt"
	"math/rand"
	"tsp-rs/internal/model"
)

// Almacena la información completa generada tras la ejecución
// de la heurística para una semilla específica
type ResultadoCorrida struct {
	Semilla      int64
	Trayectoria  []int
	Costo        float64
	Estadisticas Estadisticas
	EsFactible   bool
	Historial    []PuntoConvergencia
	Err          error
}

/* Registra una muestra en el tiempo (número de evaluación y costo actual)
 * para graficar la curva de convergencia de la heurística
 */
type PuntoConvergencia struct {
	Evaluacion int     `json:"evaluacion"`
	Costo      float64 `json:"costo"`
}

/* Indica que la búsqueda en el umbral actual se estancó y no fue
 * posible aceptar la cantidad requerida de soluciones dentro de la cota de intentos
 */
var ErrLoteIncompleto = errors.New("no se pudo completar el lote dentro del límite de intentos")

/*Explora el vecindario del estado actual 's' bajo el umbral de aceptación 'T'
 * genera vecinos iterativamente hasta completar la cuota 'params.Lote' o alcanzar 'params.MaxIntentosPorLote'
 * regresa el promedio de costos de las soluciones aceptadas en el lote, el último estado aceptado y un error
 */
func CalculaLote(rng *rand.Rand, grafica *model.GraficaTSP, T float64, s []int, N float64, params Parametros, stats *Estadisticas) (float64, []int, error) {
	fS, err := f(grafica, s, N)

	if err != nil {
		return 0, nil, err
	}

	c := 0
	var r float64
	intentos := 0

	// Genera soluciones vecinas hasta acumular la cantidad deseada de aceptaciones
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

		// Criterio de Aceptación por Umbrales (TA): acepta si empeora como máximo en 'T'
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

/* Realiza un muestreo de 'iteraciones' pasos para estimar
 * la proporción empírica de movimientos aceptados a un umbral de temperatura 'T' determinado.
 */
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

/* ResolverTSP coordina el flujo completo para solucionar el TSP dado una semilla:
 *1. inicialización del PRNG y verificación del normalizador N
 *2. generación de solución inicial mediante Vecino Más Cercano
 *3. calibración automática de la temperatura inicial (T0)
 *4. ejecución del algoritmo Aceptación por Umbrales
 *5. ajuste fino final mediante Búsqueda Local 2-opt
 */
func ResolverTSP(semilla int64, grafica *model.GraficaTSP, params Parametros) ([]int, float64, Estadisticas, bool, error) {
	// Inicializa un generador pseudoaleatorio independiente para garantizar determinismo por semilla
	rng := rand.New(rand.NewSource(semilla))

	N := model.Normalizador(grafica)

	if N == 0 {
		return nil, 0, Estadisticas{}, false, fmt.Errorf("El normalizador N es 0: porfavor revisa la gráfica de entrada")
	}

	// Construye la solución inicial con heurística voraz
	s := SolucionInicialVecinoCercano(rng, grafica, N)

	// Calibra la temperatura inicial según la tasa de aceptación objetivo
	T0, err := TemperaturaInicial(rng, grafica, s, 8.0, params.Aceptacion, params.EpsilonP, N, params.IteracionesPorcentaje)

	if err != nil {
		return nil, 0, Estadisticas{}, false, fmt.Errorf("Calculando temperatura inicial: %w", err)
	}

	// Ejecuta el proceso metaheurístico principal
	mejorS, mejorCosto, estadisticas, err := AceptacionPorUmbrales(rng, grafica, T0, s, N, params)
	if err != nil {
		return nil, 0, estadisticas, false, fmt.Errorf("Corriendo aceptación por umbrales: %w", err)
	}

	// Barrido final con búsqueda local determinista (Hill-Climbing) para refinar la solución obtenida
	mejorS, mejorCosto, err = BusquedaLocal(grafica, mejorS, N)
	if err != nil {
		return nil, 0, estadisticas, false, fmt.Errorf("corriendo barrido final: %w", err)
	}

	mejorEsFactible := model.EsFactible(grafica, mejorS)
	return mejorS, mejorCosto, estadisticas, mejorEsFactible, nil
}

/* Ejecuta una optimización determinista aplicando el operador
 * 2-opt sobre todos los pares de aristas hasta alcanzar un óptimo local donde ningún movimiento mejore el costo
 */
func BusquedaLocal(grafica *model.GraficaTSP, s []int, N float64) ([]int, float64, error) {
	actual := append([]int(nil), s...)
	costoActual, err := f(grafica, actual, N)
	if err != nil {
		return nil, 0, err
	}

	n := len(actual)
	mejorando := true

	// Explora exhaustivamente el vecindario mientras exista una mejora estricta
	for mejorando {
		mejorando = false

		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				candidato := aplicar2opt(actual, i, j)
				costoCandidato, err := f(grafica, candidato, N)
				if err != nil {
					return nil, 0, err
				}
				// Criterio de aceptación voraz
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

// Invierte el orden del subsegmento comprendido entre los índices i y j en la solución s
func aplicar2opt(s []int, i, j int) []int {
	candidato := make([]int, len(s))
	copy(candidato, s)
	for a, b := i, j; a < b; a, b = a+1, b-1 {
		candidato[a], candidato[b] = candidato[b], candidato[a]
	}
	return candidato
}
