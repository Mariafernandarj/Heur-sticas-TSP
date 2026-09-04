package main

import (
	"log"
	"tsp-rs/internal/data"
	"tsp-rs/internal/model"
)

func main() {
	// Se recibe el archivo .tsp
	ids, err := data.RecibirArchivo()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[1/4] OK: Archivo procesado correctamente (%d IDs obtenidos)\n", len(ids))

	log.Println("[2/4] Conectando a la base de datos 'tsp.sql'...")

	// SE conecta la base de datos
	db, err := data.InicioDB("tsp.sql")
	if err != nil {
		log.Fatal("Error iniciando la BD:", err)
	}
	defer db.Close()
	log.Println("[2/4] OK: Conexión establecida a la base de datos.")

	log.Println("[3/4] Construyendo matriz de adyacencias...")
	// Se construye la grafica
	grafica, err := model.ConstruirMatrizAdyacencias(db, ids)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[3/4] OK: Gráfica construida en memoria.")

	log.Println("[4/4] Imprimiendo gráfica resultante:")
	model.ImprimirGrafica(grafica)

	log.Println("[4/4] Imprimiendo gráfica resultante:")
}
