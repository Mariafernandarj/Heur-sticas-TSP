package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"tsp-rs/internal/data"
	"tsp-rs/internal/model"
	"tsp-rs/internal/reporte"
	"tsp-rs/internal/rs"
)

func main() {
	multiPtr := flag.Bool("multi", false, "Si es true, corre varias semillas en paralelo y se queda con la mejor")
	semillaPtr := flag.Int64("semilla", 42, "Semilla para una ejecución simple")
	semillasPtr := flag.String("semillas", "", "Semillas para modo multi, separadas por comas")
	pathPtr := flag.String("path", "input.tsp", "Ruta del archivo de entrada .tsp")

	flag.Parse()

	ids, err := data.RecibirArchivo(*pathPtr)

	if err != nil {
		log.Fatal(err)
	}

	db, err := data.InicioDB("tsp.sql")
	if err != nil {
		log.Fatal("Error iniciando la BD:", err)
	}
	defer db.Close()

	grafica, err := model.ConstruirMatrizAdyacencias(db, ids)
	if err != nil {
		log.Fatal(err)
	}

	peso, arista, err := model.DistanciaMaxima(*grafica)
	if err != nil {
		log.Fatal("Error calculando distancia máxima:", err)
	}
	log.Printf("Distancia máxima encontrada: %.2f (entre %v y %v)\n", peso, arista.U, arista.V)

	// --- Prueba de Normalizador ---
	N := model.Normalizador(grafica)
	log.Printf("Valor de normalización (N): %.2f\n", N)

	// --- Aceptación por umbrales ---
	log.Println("Corriendo la heurística")

	params, err := rs.LeerParametros("config.json")
	if err != nil {
		log.Fatal("Error leyendo los parámetros:", err)
	}

	if *multiPtr {
		if *semillasPtr == "" {
			log.Fatal("En modo -multi debes proporcionar -seeds")
		}

		semillas, err := parsearSemillas(*semillasPtr)
		if err != nil {
			log.Fatal(err)
		}
		correrEnParalelo(grafica, semillas, params)
		return
	}

	inicio := time.Now()
	trayectoria, costo, estadisticas, esFactible, err := rs.ResolverTSP(*semillaPtr, grafica, params)

	if err != nil {
		log.Fatal("Error corriendo aceptación por umbrales:", err)
	}
	duracion := time.Since(inicio)

	log.Printf("OK: Heurística terminada en %s\n", duracion)

	log.Println("Resultado:")
	log.Printf("  Costo (normalizado): %.6f\n", costo)
	log.Printf("  Trayectoria (%d ciudades): %v\n", len(trayectoria), trayectoria)
	log.Printf("  Solución final factible: %v\n", esFactible)
	log.Printf("  Soluciones aceptadas durante la corrida: %d (factibles: %d [%.1f%%], no factibles: %d [%.1f%%])\n",
		estadisticas.SolucionesAceptadas,
		estadisticas.SolucionesFactibles, 100*estadisticas.PorcentajeFactibles(),
		estadisticas.SolucionesNoFactibles, 100*(1-estadisticas.PorcentajeFactibles()))

	// Guardar reporte
	resultadoReporte := reporte.ResultadoReporte{
		Semilla:     *semillaPtr,
		Costo:       costo,
		Factible:    esFactible,
		Trayectoria: trayectoria,
		Historial:   estadisticas.Historial,
	}

	rep := reporte.Reporte{
		Fecha:       time.Now(),
		Modo:        "simple",
		Semillas:    []int64{*semillaPtr},
		Parametros:  params,
		TiempoTotal: duracion,
		Resultados:  []reporte.ResultadoReporte{resultadoReporte},
	}

	err = reporte.Guardar(rep, "resultados.json")
	if err != nil {
		log.Printf("Error guardando reporte: %v", err)
	}

	err = reporte.GuardarTXT(rep, "resultados")
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

func correrEnParalelo(grafica *model.GraficaTSP, semillas []int64,
	params rs.Parametros) {
	log.Printf(" Modo -multi: corriendo %d semillas en paralelo...\n", len(semillas))

	//params, err := rs.LeerParametros("config.json")

	inicio := time.Now()
	resultados := rs.ResolverTspParalelo(grafica, params, semillas)
	duracion := time.Since(inicio)

	for _, r := range resultados {
		if r.Err != nil {
			log.Printf("  semilla=%d -> ERROR: %v\n", r.Semilla, r.Err)
			continue
		}
		log.Printf("  semilla=%d -> costo=%.6f, factible=%v (aceptadas: %d, factibles: %d [%.1f%%], no factibles: %d [%.1f%%])\n",
			r.Semilla, r.Costo, r.EsFactible,
			r.Estadisticas.SolucionesAceptadas,
			r.Estadisticas.SolucionesFactibles, 100*r.Estadisticas.PorcentajeFactibles(),
			r.Estadisticas.SolucionesNoFactibles, 100*(1-r.Estadisticas.PorcentajeFactibles()),
		)
	}

	mejor, ok := rs.MejorCorrida(resultados)
	if !ok {
		log.Fatal("Todas las corridas fallaron")
	}

	log.Printf("Mejor resultado: semilla=%d, costo=%.6f, factible=%v (tardó %s en total)\n",
		mejor.Semilla, mejor.Costo, mejor.EsFactible, duracion)
	log.Printf("  Trayectoria (%d ciudades): %v\n", len(mejor.Trayectoria), mejor.Trayectoria)

	// Preparar resultados del reporte
	resultadosReporte := make([]reporte.ResultadoReporte, 0, len(resultados))

	for _, r := range resultados {
		if r.Err != nil {
			continue
		}

		resultadosReporte = append(resultadosReporte, reporte.ResultadoReporte{
			Semilla:     r.Semilla,
			Costo:       r.Costo,
			Factible:    r.EsFactible,
			Trayectoria: r.Trayectoria,
			Historial:   r.Estadisticas.Historial,
		})
	}

	// Guardar reporte
	rep := reporte.Reporte{
		Fecha:       time.Now(),
		Modo:        "multi",
		Semillas:    semillas,
		Parametros:  params,
		TiempoTotal: duracion,
		Resultados:  resultadosReporte,
	}

	if err := reporte.Guardar(rep, "reporte_multi.json"); err != nil {
		log.Printf("Error guardando reporte: %v", err)
	} else {
		log.Println("Reporte guardado correctamente.")
	}

	if err := reporte.Guardar(rep, "reporte_multi.json"); err != nil {
		log.Printf("Error guardando reporte: %v", err)
	} else {
		log.Println("Reporte guardado correctamente.")
	}

	carpeta := "graficas"

	rutaDat, rutaScript, err := reporte.GuardarConvergenciaGnuplot(carpeta, rep.Modo, rep.Resultados)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Datos de convergencia: %s\n", rutaDat)
	log.Printf("Script de gnuplot: %s\n", rutaScript)
}

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
