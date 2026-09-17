package rs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parametros struct {
	Lote                   int     `json:"lote"`                      // Tamaño del lote
	MaxIntentosPorLote     int     `json:"max_intentos_por_lote"`     //Cota de intentos para completar un lote
	MaxLotesPorTemperatura int     `json:"max_lotes_por_temperatura"` //Cota de lotes seguidos "sin empeorar" antes de forzar el enfriamiento
	Phi                    float64 `json:"phi"`                       // Factor de enfriamineto
	Epsilon                float64 `json:"epsilon"`
	Aceptacion             float64 `json:"aceptacion"` // Porcentaje de aceptación objetivo para T
	EpsilonP               float64 `json:"epsilon_p"`
	IteracionesPorcentaje  int     `json:"iteraciones_porcentaje"` //N en porcentajes de aceptado
}

func LeerParametros(ruta string) (Parametros, error) {
	archivo, err := os.ReadFile(ruta)
	if err != nil {
		return Parametros{}, fmt.Errorf("No se pudo leer el archivo de parámetros: %w", err)
	}

	var parametros Parametros

	err = json.Unmarshal(archivo, &parametros)
	if err != nil {
		return Parametros{}, fmt.Errorf("Error al interpretar el archivo de parámetros: %w", err)
	}
	return parametros, nil
}
