package main

import (
	"log"
	"time"
	"tsp-rs/internal/data"
	"tsp-rs/internal/model"
	"tsp-rs/internal/rs"
)

func main() {
	// Se recibe el archivo .tsp
	ids, err := data.RecibirArchivo()

	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("[1/4] OK: Archivo procesado correctamente (%d IDs obtenidos)\n", len(ids))

	//log.Println("[2/4] Conectando a la base de datos 'tsp.sql'...")

	// SE conecta la base de datos
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
	//log.Println("[3/4] Construyendo matriz de adyacencias...")
	// Se construye la grafica
	grafica, err := model.ConstruirMatrizAdyacencias(db, ids)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[3/4] OK: Gráfica construida en memoria.")

	// PRUEBA: gráfica original

	//log.Println("Matriz de adyacencias original:")
	//model.ImprimirGrafica(grafica)

	// --- Prueba de DistanciaMaxima ---
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

}
