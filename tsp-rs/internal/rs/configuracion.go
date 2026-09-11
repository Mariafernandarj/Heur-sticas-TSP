package rs

import ()

type Parametros struct {
	Lote                   int     // Tamaño del lote
	MaxIntentosPorLote     int     //Cota de intentos para completar un lote
	MaxLotesPorTemperatura int     //Cota de lotes seguidos "sin empeorar" antes de forzar el enfriamiento
	Phi                    float64 // Factor de enfriamineto
	Epsilon                float64
	Aceptacion             float64 // Porcentaje de aceptación objetivo para T
	EpsilonP               float64
	IteracionesPorcentaje  int //N en porcentajes de aceptado
}

func ParametrosPorDefecto() Parametros {
	return Parametros{
		Lote:                   200,  //50,
		MaxIntentosPorLote:     5000, //2000,
		MaxLotesPorTemperatura: 200,
		Phi:                    0.99, //0.95,
		Epsilon:                1e-6,
		Aceptacion:             0.90,
		EpsilonP:               0.02,
		IteracionesPorcentaje:  200,
	}
}
