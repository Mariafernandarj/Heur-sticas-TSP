package rs

import (
	"errors"
	"math"
	"math/rand"
	"tsp-rs/internal/model"
)

/* Ejecuta la metaheurística de Aceptación por Umbrales (TA) para resolver el TSP
 * recibe un generador aleatorio 'rng', el grafo, el umbral inicial 'T', la solución inicial 's',
 * el factor de normalización 'N' y la estructura de parámetros de control
 */
func AceptacionPorUmbrales(rng *rand.Rand, grafica *model.GraficaTSP, T float64, s []int, N float64, params Parametros) ([]int, float64, Estadisticas, error) {

	estadisticas := Estadisticas{}

	// Mantiene un registro independiente de la mejor solución y costo globales
	mejorS := append([]int(nil), s...)
	mejorCosto, err := f(grafica, s, N)

	if err != nil {
		return nil, 0, estadisticas, err
	}

	// Registra la solución inicial en el historial de convergencia
	estadisticas.Historial = append(estadisticas.Historial, PuntoConvergencia{Evaluacion: 0, Costo: mejorCosto})

	p := 0.0

	// Bucle principal del esquema de enfriamiento (mientras el umbral supere epsilon)
	for T > params.Epsilon {
		q := math.Inf(1)
		lotes := 0
		// Bucle de equilibrio a umbral constante: se ejecuta mientras el desempeño no empeore respecto al lote anterior
		for p <= q {
			if lotes >= params.MaxLotesPorTemperatura {
				estadisticas.TopesPorMaxLotes++
				break
			}
			lotes++
			q = p

			// Procesa un nuevo lote de vecino candidato dentro del umbral T actual
			pNuevo, sNuevo, err := CalculaLote(rng, grafica, T, s, N, params, &estadisticas)
			if err != nil {
				// Si el lote se interrumpe por falta de vecinos aceptados, el umbral es demasiado estrecho
				// Se asume convergencia local y se regresa de forma segura la mejor solución encontrada hasta el momento
				if errors.Is(err, ErrLoteIncompleto) {
					return mejorS, mejorCosto, estadisticas, nil
				}
				return nil, 0, estadisticas, err
			}
			p, s = pNuevo, sNuevo

			// Evalúa la solución actual del lote y actualiza la mejor global si se encuentra una mejor
			costoActual, err := f(grafica, s, N)
			if err != nil {
				return nil, 0, estadisticas, err
			}
			if costoActual < mejorCosto {
				mejorCosto = costoActual
				mejorS = append([]int(nil), s...)
				estadisticas.Historial = append(estadisticas.Historial, PuntoConvergencia{
					Evaluacion: estadisticas.EvaluacionesTotales,
					Costo:      mejorCosto,
				})
			}
		}
		// Reduce el umbral utilizando el coeficiente de enfriamiento phi
		T *= params.Phi
	}
	return mejorS, mejorCosto, estadisticas, nil
}

/* Realiza una búsqueda en el intervalo [T1, T2] para hallar un umbral de temperatura
 * cuyo porcentaje de vecinos aceptados esté dentro de un margen 'epsilonP' respecto a la 'Aceptacion' objetivo
 */
func BusquedaBinaria(rng *rand.Rand, grafica *model.GraficaTSP, s []int, T1, T2, Aceptacion, epsilonP, N float64, iteraciones int) (float64, error) {
	Tm := (T1 + T2) / 2

	// Criterio de parada por amplitud del intervalo
	if T2-T1 < epsilonP {
		return Tm, nil
	}

	// Calcula el porcentaje real de aceptados para la temperatura media Tm
	p, err := PorcentajeAceptados(rng, grafica, s, Tm, N, iteraciones)

	if err != nil {
		return 0, err
	}

	// Criterio de parada por tolerancia alcanzada en la tasa de aceptación
	if math.Abs(Aceptacion-p) < epsilonP {
		return Tm, nil
	}

	// Ajuste del intervalo según la tasa de aceptación obtenida
	if p > Aceptacion {
		return BusquedaBinaria(rng, grafica, s, T1, Tm, Aceptacion, epsilonP, N, iteraciones)
	}

	return BusquedaBinaria(rng, grafica, s, Tm, T2, Aceptacion, epsilonP, N, iteraciones)
}

/* Determina la temperatura inicial óptima ajustada al porcentaje de aceptación deseado.
 * aplica una fase de acotamiento exponencial seguida de búsqueda binaria
 */
func TemperaturaInicial(rng *rand.Rand, grafica *model.GraficaTSP, s []int, T, Aceptacion, epsilonP, N float64, iteraciones int) (float64, error) {
	// Evalúa la temperatura de prueba dada
	p, err := PorcentajeAceptados(rng, grafica, s, T, N, iteraciones)
	if err != nil {
		return 0, err
	}

	// Si la temperatura inicial ya cumple la tolerancia, se retorna directamente
	if math.Abs(Aceptacion-p) <= epsilonP {
		return T, nil
	}

	var T1, T2 float64
	// Fase de acotamiento del intervalo [T1, T2]
	if p < Aceptacion {
		// Aumenta T exponencialmente si la tasa de aceptación es demasiado baja
		for p < Aceptacion {
			T *= 2
			p, err = PorcentajeAceptados(rng, grafica, s, T, N, iteraciones)
			if err != nil {
				return 0, err
			}
		}
		T1, T2 = T/2, T
	} else {
		// Reduce T exponencialmente si la tasa de aceptación es demasiado alta
		for p > Aceptacion {
			T /= 2
			p, err = PorcentajeAceptados(rng, grafica, s, T, N, iteraciones)
			if err != nil {
				return 0, err
			}
		}
		T1, T2 = T, T*2
	}
	// Refina el rango hallado usando Búsqueda Binaria
	return BusquedaBinaria(rng, grafica, s, T1, T2, Aceptacion, epsilonP, N, iteraciones)
}
