package data

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*Lee un archivo de texto con números separados por comas
 *(como input-n.tsp) y devuelve los datos procesados como un slice de enteros
 */
func CargarArchivo(path string) ([]int, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contenido := strings.TrimSpace(string(bytes))
	partes := strings.Split(contenido, ",")

	datos := make([]int, 0, len(partes))
	for _, p := range partes {
		p = strings.TrimSpace(p)
		if val, err := strconv.Atoi(p); err == nil {
			datos = append(datos, val)
		}
	}
	return datos, nil
}

/* Procesa la carga de un archivo dada su ruta e informa
 * por consola la cantidad de elementos leídos correctamente
 */
func RecibirArchivo(ruta string) ([]int, error) {
	datos, err := CargarArchivo(ruta)
	if err != nil {
		return nil, fmt.Errorf("al cargar archivo: %w", err)
	}
	fmt.Printf("Se leyeron ls %d elementos correctamente:\n", len(datos))

	return datos, nil
}
