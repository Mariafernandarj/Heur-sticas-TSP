package rs

import (
	"encoding/json"
	"fmt"
	"os"
)

/* Define los hiperparámetros de configuración requeridos
 * para la ejecución y el ajuste del algoritmo de Recocido Simulado
 */
type Parametros struct {
	Lote                   int     `json:"lote"`                      // Cantidad de soluciones aceptadas requeridas por lote
	MaxIntentosPorLote     int     `json:"max_intentos_por_lote"`     // Cota máxima de intentos permitidos para completar un lote
	MaxLotesPorTemperatura int     `json:"max_lotes_por_temperatura"` // Lotes consecutivos sin mejora antes de reducir la temperatura
	Phi                    float64 `json:"phi"`                       // Coeficiente de enfriamiento (factor de reducción de T)
	Epsilon                float64 `json:"epsilon"`                   // Umbral de temperatura mínima (criterio de parada)
	Aceptacion             float64 `json:"aceptacion"`                // Tasa de aceptación objetivo para calibrar la temperatura inicial
	EpsilonP               float64 `json:"epsilon_p"`                 // Margen de error en la tasa de aceptación
	IteracionesPorcentaje  int     `json:"iteraciones_porcentaje"`    // Cantidad de soluciones aceptadas requeridas por lote
}

/* Lee un archivo en formato JSON desde la ruta especificada
 * y deserializa sus datos en una estructura de tipo Parametros
 */
func LeerParametros(ruta string) (Parametros, error) {
	// Lee la totalidad del archivo de configuración desde el sistema de archivos
	archivo, err := os.ReadFile(ruta)
	if err != nil {
		return Parametros{}, fmt.Errorf("No se pudo leer el archivo de parámetros: %w", err)
	}

	var parametros Parametros

	// Deserializa la estructura JSON en la variable de parámetros
	err = json.Unmarshal(archivo, &parametros)
	if err != nil {
		return Parametros{}, fmt.Errorf("Error al interpretar el archivo de parámetros: %w", err)
	}
	return parametros, nil
}
