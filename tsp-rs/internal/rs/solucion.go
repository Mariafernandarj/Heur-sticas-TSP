package rs

import (
	"math"
	"math/rand"
	"tsp-rs/internal/model"
)

// Estadisticas recopila métricas de rendimiento y monitoreo durante la ejecución de la heurística
type Estadisticas struct {
	SolucionesAceptadas   int                 // Cantidad total de soluciones que cumplieron el criterio de aceptación
	SolucionesFactibles   int                 // Cantidad de soluciones aceptadas que constituyen un tour válido
	SolucionesNoFactibles int                 // Cantidad de soluciones aceptadas que utilizan aristas penalizadas
	TopesPorMaxLotes      int                 // Contador de veces que se forzó el enfriamiento al alcanzar la cota de lotes
	EvaluacionesTotales   int                 // Número acumulado de evaluaciones de la función objetivo
	Historial             []PuntoConvergencia // Registro histórico de mejoras para trazado de convergencia
}

// Genera un recorrido inicial (permutación) de forma completamente aleatoria
func SolucionInicial(rng *rand.Rand, grafica *model.GraficaTSP) []int {
	n := len(grafica.Ciudades)
	indices := rng.Perm(n)
	s := make([]int, n)
	for i, idx := range indices {
		s[i] = idx
	}
	return s
}

/* Construye una solución inicial mediante la heurística voraz (greedy) del Vecino Más Cercano
 * elecciona un nodo inicial de forma aleatoria y conecta iterativamente la ciudad más próxima aún no visitada
 */
func SolucionInicialVecinoCercano(rng *rand.Rand, grafica *model.GraficaTSP, N float64) []int {
	n := len(grafica.Ciudades)
	visitado := make([]bool, n)

	// Selecciona un punto de partida al azar para diversificar las semillas iniciales
	actual := rng.Intn(n)
	visitado[actual] = true

	path := make([]int, 0, n)
	path = append(path, actual)

	// Construye el recorrido seleccionando siempre el vecino de menor costo
	for len(path) < n {
		mejorIdx := -1
		mejorPeso := math.Inf(1)

		for j := 0; j < n; j++ {
			if visitado[j] {
				continue
			}
			peso := model.PesoAumentado(grafica, actual, j, N)
			if peso < mejorPeso {
				mejorPeso = peso
				mejorIdx = j
			}
		}

		visitado[mejorIdx] = true
		path = append(path, mejorIdx)
		actual = mejorIdx
	}

	return path
}

/* Calcula la función objetivo evaluando la trayectoria de una solución candidata 's'
 * y retornando su costo total normalizado
 */
func f(grafica *model.GraficaTSP, s []int, N float64) (float64, error) {
	evaluacion, err := model.EvaluarTrayectoria(grafica, s, N)
	if err != nil {
		return 0, err
	}
	return model.Costo(evaluacion, N), nil
}

/* Calcula la proporción de soluciones válidas (factibles)
 * con respecto al total de soluciones aceptadas durante la ejecución
 */
func (e Estadisticas) PorcentajeFactibles() float64 {
	if e.SolucionesAceptadas == 0 {
		return 0
	}
	return float64(e.SolucionesFactibles) / float64(e.SolucionesAceptadas)
}
