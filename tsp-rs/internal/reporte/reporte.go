package reporte

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tsp-rs/internal/rs"
)

type Reporte struct {
	Fecha       time.Time          `json:"fecha"`
	Modo        string             `json:"modo"`
	Semillas    []int64            `json:"semillas"`
	TiempoTotal time.Duration      `json:"tiempo_total"`
	Parametros  rs.Parametros      `json:"parametros"`
	Resultados  []ResultadoReporte `json:"resultados"`
}

type ResultadoReporte struct {
	Semilla     int64   `json:"semilla"`
	Evaluacion  float64 `json:"evaluacion"`
	Costo       float64 `json:"costo"`
	Factible    bool    `json:"factible"`
	Trayectoria []int   `json:"trayectoria"`
}

func Guardar(reporte Reporte, carpeta string) error {

	if err := os.MkdirAll(carpeta, 0755); err != nil {
		return fmt.Errorf("creando carpeta de reportes: %w", err)
	}

	nombre := fmt.Sprintf(
		"ejecucion_%s.json",
		time.Now().Format("2006-01-02_15-04-05"),
	)

	ruta := filepath.Join(carpeta, nombre)

	datos, err := json.MarshalIndent(
		reporte,
		"",
		"  ",
	)

	if err != nil {
		return fmt.Errorf("serializando reporte: %w", err)
	}

	if err := os.WriteFile(
		ruta,
		datos,
		0644,
	); err != nil {
		return fmt.Errorf("guardando reporte: %w", err)
	}

	return nil
}

func GuardarTXT(reporte Reporte, carpeta string) error {
	if err := os.MkdirAll(carpeta, 0755); err != nil {
		return fmt.Errorf("creando carpeta de reportes: %w", err)
	}
	nombre := fmt.Sprintf("ejecucion_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	ruta := filepath.Join(carpeta, nombre)
	archivo, err := os.Create(ruta)

	if err != nil {
		return fmt.Errorf("creando reporte TXT: %w", err)
	}
	defer archivo.Close()
	fmt.Fprintln(archivo, "========================================")
	fmt.Fprintln(archivo, " REPORTE DE EJECUCIÓN TSP")
	fmt.Fprintln(archivo, "========================================")
	fmt.Fprintln(archivo)

	fmt.Fprintf(archivo, "Fecha: %s\n", reporte.Fecha.Format("02/01/2006 15:04:05"))
	fmt.Fprintf(archivo, "Modo: %s\n", reporte.Modo)
	fmt.Fprintf(archivo, "Semillas: %v\n", reporte.Semillas)
	fmt.Fprintf(archivo, "Tiempo total: %s\n", reporte.TiempoTotal)

	fmt.Fprintln(archivo)
	fmt.Fprintln(archivo, "----------------------------------------")
	fmt.Fprintln(archivo, "PARÁMETROS")
	fmt.Fprintln(archivo, "----------------------------------------")
	fmt.Fprintln(archivo)
	fmt.Fprintf(archivo, "%+v\n", reporte.Parametros)

	fmt.Fprintln(archivo)
	fmt.Fprintln(archivo, "----------------------------------------")
	fmt.Fprintln(archivo, "RESULTADOS")
	fmt.Fprintln(archivo, "----------------------------------------")

	for _, resultado := range reporte.Resultados {
		fmt.Fprintln(archivo)
		fmt.Fprintf(archivo, "Semilla: %d\n", resultado.Semilla)
		fmt.Fprintf(archivo, "Evaluación: %.6f\n", resultado.Evaluacion)
		fmt.Fprintf(archivo, "Costo: %.6f\n", resultado.Costo)
		fmt.Fprintf(archivo, "Factible: %v\n", resultado.Factible)
		fmt.Fprintf(archivo, "Trayectoria: %v\n", resultado.Trayectoria)
	}
	fmt.Fprintln(archivo)
	fmt.Fprintln(archivo, "========================================")
	return nil
}

func ConvertirResultado(r rs.ResultadoCorrida) ResultadoReporte {
	return ResultadoReporte{
		Semilla:     r.Semilla,
		Evaluacion:  r.Evaluacion,
		Costo:       r.Costo,
		Factible:    r.EsFactible,
		Trayectoria: r.Trayectoria,
	}
}
