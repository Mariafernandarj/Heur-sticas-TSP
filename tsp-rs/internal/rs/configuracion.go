package rs

import (
	"fmt"
)

type Parametros struct {
	Lote                  int
	MaxIntentosLote       int
	Phi                   float64
	Epsilon               float64
	Aceptacion            float64
	EpsilonP              float64
	IteracionesPorcentaje int
}

func ParametrosPorDefecto() Parametros {
	return Parametros{
		Lote:                  50,
		MaxIntentosLote:       2000,
		Phi:                   0.95,
		Epsilon:               1e-6,
		Aceptacion:            0.90,
		EpsilonP:              0.02,
		IteracionesPorcentaje: 200,
	}
}
