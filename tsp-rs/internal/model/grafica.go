package model

import (
	"database/sql"
	"fmt"
	"math"
	"tsp-rs/internal/data"
)

type GraficaTSP struct {
	Ciudades    []data.Ciudad
	IndicePorID map[int]int
	Matriz      [][]float64
}

func ConstruirMatrizAdyacencias(db *sql.DB, ids []int) (*GraficaTSP, error) {
	ciudades, indicePorID, err := data.GetCiudades(db, ids)
	if err != nil {
		return nil, err
	}

	matriz := inicializarMatriz(len(ciudades))

	if err := cargarConexiones(db, ciudades, indicePorID, matriz); err != nil {
		return nil, err
	}

	return &GraficaTSP{
		Ciudades:    ciudades,
		IndicePorID: indicePorID,
		Matriz:      matriz,
	}, nil
}

func inicializarMatriz(n int) [][]float64 {
	matriz := make([][]float64, n)

	for i := range matriz {
		matriz[i] = make([]float64, n)

		for j := range matriz[i] {
			if i == j {
				matriz[i][j] = 0
			} else {
				matriz[i][j] = math.Inf(1)
			}
		}
	}
	return matriz
}

func cargarConexiones(db *sql.DB, ciudades []data.Ciudad, indicePorID map[int]int, matriz [][]float64) error {
	filasConn, err := db.Query("SELECT id_city_1, id_city_2 FROM connections")

	if err != nil {
		return fmt.Errorf("Error consultando connections: %w", err)
	}
	defer filasConn.Close()

	for filasConn.Next() {
		var id1, id2 int

		if err := filasConn.Scan(&id1, &id2); err != nil {
			return fmt.Errorf("Error leyendo conections: %w", err)
		}
		i, okI := indicePorID[id1]
		j, okJ := indicePorID[id2]

		if !okI || !okJ {
			continue
		}

		distancia := calcularDistancia(ciudades[i], ciudades[j])

		matriz[i][j] = distancia
		matriz[i][j] = distancia

	}

	if err := filasConn.Err(); err != nil {
		return err
	}
	return nil
}

func calcularDistancia(c1, c2 data.Ciudad) float64 {
	return DistanciaNatural(c1.Latitud, c1.Longitud, c2.Latitud, c2.Longitud)
}

func ImprimirGrafica(grafica *GraficaTSP) {
	fmt.Println("\n========== GRAFO TSP ==========")

	// Mostrar ciudades
	fmt.Println("\n--- Ciudades ---")

	for i, ciudad := range grafica.Ciudades {
		fmt.Printf(
			"[%d] ID: %d | Nombre: %s | Lat: %.6f | Lon: %.6f\n",
			i,
			ciudad.ID,
			ciudad.Nombre,
			ciudad.Latitud,
			ciudad.Longitud,
		)
	}

	// Mostrar matriz
	fmt.Println("\n--- Matriz de adyacencias ---")

	// Encabezado
	fmt.Printf("%8s", "")
	for i := range grafica.Matriz {
		fmt.Printf("%12d", i)
	}
	fmt.Println()

	// Filas
	for i := range grafica.Matriz {
		fmt.Printf("%8d", i)

		for j := range grafica.Matriz[i] {
			distancia := grafica.Matriz[i][j]

			if math.IsInf(distancia, 1) {
				fmt.Printf("%12s", "INF")
			} else {
				fmt.Printf("%12.2f", distancia)
			}
		}

		fmt.Println()
	}

	fmt.Println("\n===============================")
}
