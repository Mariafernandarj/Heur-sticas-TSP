package main

import (
	"flag"
	"log"
	"time"
	"tsp-rs/internal/data"
	"tsp-rs/internal/model"
	"tsp-rs/internal/reporte"
	"tsp-rs/internal/rs"
)

func main() {
	multiPtr := flag.Bool("multi", false, "Si es true, corre varias semillas en paralelo y se queda con la mejor")

	ids, err := data.RecibirArchivo()

	if err != nil {
		log.Fatal(err)
	}

	db, err := data.InicioDB("tsp.sql")
	if err != nil {
		log.Fatal("Error iniciando la BD:", err)
	}
	defer db.Close()
	//log.Println("[2/4] OK: Conexión establecida a la base de datos.")

	totalCiudades, totalConexiones, err := data.ContarCiudadesYConexiones(db)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Ciudades: %d", totalCiudades)
	log.Printf("Conexiones: %d", totalConexiones)

	grafica, err := model.ConstruirMatrizAdyacencias(db, ids)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[3/4] OK: Gráfica construida en memoria.")

	peso, arista, err := model.DistanciaMaxima(*grafica)
	if err != nil {
		log.Fatal("Error calculando distancia máxima:", err)
	}
	log.Printf("Distancia máxima encontrada: %.2f (entre %v y %v)\n", peso, arista.U, arista.V)

	// --- Prueba de Normalizador ---
	N := model.Normalizador(grafica)
	log.Printf("Valor de normalización (N): %.2f\n", N)

	// --- Aceptación por umbrales ---
	log.Println("Corriendo la heurística de aceptación por umbrales...")
	params := rs.ParametrosPorDefecto()

	// Semilla fija
	const semilla = 42

	if *multiPtr {
		correrEnParalelo(grafica)
		return
	}

	inicio := time.Now()
	trayectoria, costo, estadisticas, esFactible, err := rs.ResolverTSP(semilla, grafica, params)

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
		Semilla:     semilla,
		Costo:       costo,
		Factible:    esFactible,
		Trayectoria: trayectoria,
	}

	rep := reporte.Reporte{
		Modo:        "simple",
		Semillas:    []int64{semilla},
		Parametros:  params,
		TiempoTotal: duracion,
		Resultados:  []reporte.ResultadoReporte{resultadoReporte},
	}
	err = reporte.Guardar(rep, "reporte_simple.json")

	if err != nil {
		log.Printf("Error guardando reporte: %v", err)
	}

}

func correrEnParalelo(grafica *model.GraficaTSP) {
	semillas := []int64{1, 2, 3, 4, 5, 6, 7, 8}
	log.Printf("[4/5] Modo -multi: corriendo %d semillas en paralelo...\n", len(semillas))

	params := rs.ParametrosPorDefecto()

	inicio := time.Now()
	resultados := rs.ResolverTspParalelo(grafica, params, semillas)
	duracion := time.Since(inicio)

	for _, r := range resultados {
		if r.Err != nil {
			log.Printf("  semilla=%d -> ERROR: %v\n", r.Semilla, r.Err)
			continue
		}
		log.Printf("  semilla=%d -> costo=%.6f, factible=%v\n", r.Semilla, r.Costo, r.EsFactible)
	}

	mejor, ok := rs.MejorCorrida(resultados)
	if !ok {
		log.Fatal("Todas las corridas fallaron")
	}

	log.Printf("[5/5] Mejor resultado: semilla=%d, costo=%.6f, factible=%v (tardó %s en total)\n",
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
		})
	}

	// Guardar reporte
	rep := reporte.Reporte{
		Modo:        "multi",
		Semillas:    semillas,
		Parametros:  params,
		TiempoTotal: duracion,
		Resultados:  resultadosReporte,
	}

	err := reporte.Guardar(rep, "reporte_multi.json")
	if err != nil {
		log.Printf("Error guardando reporte: %v", err)
	} else {
		log.Println("Reporte guardado correctamente.")
	}
}
