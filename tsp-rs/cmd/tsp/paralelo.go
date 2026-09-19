package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"tsp-rs/internal/model"
	"tsp-rs/internal/reporte"
	"tsp-rs/internal/rs"
)

/*Ejecuta múltiples instancias del TSP concurrentemente usando
 *las semillas proporcionadas, procesa los resultados e imprime la mejor opción
 */
func correrEnParalelo(grafica *model.GraficaTSP, semillas []int64,
	params rs.Parametros) {
	//log.Printf(" Modo -multi: corriendo %d semillas en paralelo...\n", len(semillas))

	// Medir el tiempo total de ejecución paralela
	inicio := time.Now()
	resultados := rs.ResolverTspParalelo(grafica, params, semillas)
	duracion := time.Since(inicio)

	// Reportar el estado individual de cada corrida
	for _, r := range resultados {
		if r.Err != nil {
			log.Printf("  semilla=%d -> ERROR: %v\n", r.Semilla, r.Err)
			continue
		}
		log.Printf("  semilla=%d -> costo=%.12f, factible=%v (aceptadas: %d, factibles: %d [%.1f%%], no factibles: %d [%.1f%%])\n",
			r.Semilla, r.Costo, r.EsFactible,
			r.Estadisticas.SolucionesAceptadas,
			r.Estadisticas.SolucionesFactibles, 100*r.Estadisticas.PorcentajeFactibles(),
			r.Estadisticas.SolucionesNoFactibles, 100*(1-r.Estadisticas.PorcentajeFactibles()),
		)
	}

	// Encontrar la mejor corrida de todo el conjunto paralelo
	mejor, ok := rs.MejorCorrida(resultados)
	if !ok {
		log.Fatal("Todas las corridas fallaron")
	}

	log.Printf("Mejor resultado: semilla=%d, costo=%.12f, factible=%v (tardó %s en total)\n",
		mejor.Semilla, mejor.Costo, mejor.EsFactible, duracion)
	log.Printf("  Trayectoria (%d ciudades): %v\n", len(mejor.Trayectoria), mejor.Trayectoria)

	// Preparar los resultados estructurados para los reportes
	resultadosReporte := make([]reporte.ResultadoReporte, 0, len(resultados))

	for _, r := range resultados {
		if r.Err != nil {
			continue
		}
		resultadoReporte := reporte.ConvertirResultado(r, grafica)

		resultadosReporte = append(resultadosReporte, resultadoReporte)
	}

	// Construir la estructura general del reporte multi-semilla
	rep := reporte.Reporte{
		Fecha:       time.Now(),
		Modo:        "multi",
		Semillas:    semillas,
		Parametros:  params,
		TiempoTotal: duracion,
		Resultados:  resultadosReporte,
	}

	// Guardar reportes en formatos JSON, TXT y datos para Gnuplot
	if err := reporte.Guardar(rep, "reporte_multi.json"); err != nil {
		log.Printf("Error guardando reporte: %v", err)
	} else {
		log.Println("Reporte guardado correctamente.")
	}

	err := reporte.GuardarTXT(rep, "resultados")
	if err != nil {
		log.Printf("Error guardando reporte TXT: %v", err)
	}

	carpeta := "graficas"

	rutaDat, rutaScript, err := reporte.GuardarConvergenciaGnuplot(carpeta, rep.Modo, rep.Resultados)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Datos de convergencia: %s\n", rutaDat)
	log.Printf("Script de gnuplot: %s\n", rutaScript)
}

/* Toma una cadena de texto con números separados por comas
 * (ej: "42, 100, 204") y los convierte en un slice de enteros de 64 bits ([]int64).
 */
func parsearSemillas(texto string) ([]int64, error) {
	partes := strings.Split(texto, ",")
	semillas := make([]int64, 0, len(partes))

	for _, parte := range partes {
		semilla, err := strconv.ParseInt(strings.TrimSpace(parte), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("semilla inválida %q: %w", parte, err)
		}

		semillas = append(semillas, semilla)
	}

	return semillas, nil
}
