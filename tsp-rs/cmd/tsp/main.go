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
	//Configuración y lectura de argumentos por linea de comandos (Banderas)
	multiPtr := flag.Bool("multi", false, "Si es true, corre varias semillas en paralelo y se queda con la mejor")
	semillaPtr := flag.Int64("semilla", 42, "Semilla para una ejecución simple")
	semillasPtr := flag.String("semillas", "", "Semillas para modo multi, separadas por comas")
	pathPtr := flag.String("path", "input.tsp", "Ruta del archivo de entrada .tsp")

	flag.Parse()

	//Inicialización de datos con la base de datos proporcionada
	ids, err := data.RecibirArchivo(*pathPtr)

	if err != nil {
		log.Fatal(err)
	}

	db, err := data.InicioDB("tsp.sql")
	if err != nil {
		log.Fatal("Error iniciando la BD:", err)
	}
	defer db.Close()

	//Construcción de la matriz de adyacencias a partir de la base de datos
	grafica, err := model.ConstruirMatrizAdyacencias(db, ids)
	if err != nil {
		log.Fatal(err)
	}

	peso, arista, err := model.DistanciaMaxima(*grafica)
	if err != nil {
		log.Fatal("Error calculando distancia máxima:", err)
	}
	model.CompletarAristas(grafica, peso)

	//Calculo de la distancia máxima y normalización de la gráfica
	log.Printf("Distancia máxima encontrada: %.12f (entre %v y %v)\n", peso, arista.U, arista.V)

	N := model.Normalizador(grafica)
	log.Printf("Valor de normalización (N): %.12f\n", N)

	//Lectura de los parámetros de la heurística
	log.Println("Corriendo la heurística")

	params, err := rs.LeerParametros("config.json")
	if err != nil {
		log.Fatal("Error leyendo los parámetros:", err)
	}

	// División del modo paralelo y el modo simple
	if *multiPtr {
		if *semillasPtr == "" {
			log.Fatal("En modo -multi debes proporcionar -semillas")
		}

		semillas, err := parsearSemillas(*semillasPtr)
		if err != nil {
			log.Fatal(err)
		}
		correrEnParalelo(grafica, semillas, params)
		return
	}

	//Ejecución del modo simple es decir una semilla solamente
	inicio := time.Now()
	trayectoria, costo, estadisticas, esFactible, err := rs.ResolverTSP(*semillaPtr, grafica, params)

	if err != nil {
		log.Fatal("Error corriendo aceptación por umbrales:", err)
	}
	duracion := time.Since(inicio)

	log.Printf("OK: Heurística terminada en %s\n", duracion)

	log.Println("Resultado:")
	log.Printf("  Costo (normalizado): %.12f\n", costo)
	log.Printf("  Trayectoria (%d ciudades): %v\n", len(trayectoria), trayectoria)
	log.Printf("  Solución final factible: %v\n", esFactible)
	log.Printf("  Soluciones aceptadas durante la corrida: %d (factibles: %d [%.1f%%], no factibles: %d [%.1f%%])\n",
		estadisticas.SolucionesAceptadas,
		estadisticas.SolucionesFactibles, 100*estadisticas.PorcentajeFactibles(),
		estadisticas.SolucionesNoFactibles, 100*(1-estadisticas.PorcentajeFactibles()))

	// Generación de reportes y archivos de salida (JSON, TXT, Gnuplot)
	resultadoCorrida := rs.ResultadoCorrida{
		Semilla:      *semillaPtr,
		Costo:        costo,
		EsFactible:   esFactible,
		Trayectoria:  trayectoria,
		Estadisticas: estadisticas,
		Historial:    estadisticas.Historial,
	}

	resultadoReporte := reporte.ConvertirResultado(resultadoCorrida, grafica)

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
