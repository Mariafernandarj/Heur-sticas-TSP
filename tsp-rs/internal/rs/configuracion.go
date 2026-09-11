package rs

import ()

type Parametros struct {
	Lote                   int
	MaxIntentosPorLote     int
	MaxLotesPorTemperatura int // cota de lotes seguidos "sin empeorar" antes de forzar el enfriamiento
	Phi                    float64
	Epsilon                float64
	Aceptacion             float64
	EpsilonP               float64
	IteracionesPorcentaje  int
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
