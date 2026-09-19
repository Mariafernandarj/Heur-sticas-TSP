package data

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

type Ciudad struct {
	ID       int
	Nombre   string
	Latitud  float64
	Longitud float64
}

/* Crea una base de datos SQLite efímera en memoria y ejecuta un script
 * SQL para inicializar el esquema y los datos iniciales.
 */
func InicioDB(scriptPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	script, err := os.ReadFile(scriptPath)
	if err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(script)); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// GetCiudades recupera de la base de datos las ciudades cuyos IDs coincidan con el slice proporcionado.
// Devuelve las ciudades encontradas, un mapa con la posición de cada ID en la lista resultante
// (útil para mapear índices en matrices) y un posible error.
func GetCiudades(db *sql.DB, ids []int) ([]Ciudad, map[int]int, error) {
	filas, err := db.Query("SELECT id, name, latitude, longitude FROM cities ORDER BY id")

	if err != nil {
		return nil, nil, fmt.Errorf("Error consultando cities: %w", err)
	}
	defer filas.Close()

	idsSolicitados := make(map[int]bool)

	for _, id := range ids {
		idsSolicitados[id] = true
	}

	var ciudades []Ciudad
	indicePorID := make(map[int]int)

	for filas.Next() {
		var c Ciudad

		if err := filas.Scan(&c.ID, &c.Nombre, &c.Latitud, &c.Longitud); err != nil {
			return nil, nil, fmt.Errorf("Error leyendo ciudad: %w", err)
		}

		if !idsSolicitados[c.ID] {
			continue
		}

		indicePorID[c.ID] = len(ciudades)
		ciudades = append(ciudades, c)
	}

	if err := filas.Err(); err != nil {
		return nil, nil, err
	}

	if len(ciudades) == 0 {
		return nil, nil, fmt.Errorf("No se encontraron ciudades en la base de datos")
	}

	return ciudades, indicePorID, nil
}
