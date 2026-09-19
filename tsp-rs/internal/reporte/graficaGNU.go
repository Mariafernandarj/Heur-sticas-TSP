package reporte

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*GuardarConvergenciaGnuplot genera los archivos de datos (.dat) y de comandos (.gnuplot)
 *necesarios para graficar la curva de convergencia de cada corrida mediante Gnuplot
 */
func GuardarConvergenciaGnuplot(carpeta, modo string, resultados []ResultadoReporte) (string, string, error) {
	if err := os.MkdirAll(carpeta, 0755); err != nil {
		return "", "", fmt.Errorf("creando carpeta de reportes: %w", err)
	}

	marca := time.Now().Format("01-02_15-04-05")
	nombreDat := fmt.Sprintf("convergencia_%s.dat", marca)
	nombreScript := fmt.Sprintf("convergencia_%s.gnuplot", marca)

	rutaDat := filepath.Join(carpeta, nombreDat)
	rutaScript := filepath.Join(carpeta, nombreScript)

	// --- archivo .dat ---
	var dat strings.Builder
	for i, r := range resultados {
		if len(r.Historial) == 0 {
			continue
		}
		if i > 0 {
			dat.WriteString("\n\n") // separador de bloque (index) para gnuplot
		}
		dat.WriteString(fmt.Sprintf("# semilla %d\n", r.Semilla))
		dat.WriteString("# evaluacion\tcosto\n")
		for _, p := range r.Historial {
			dat.WriteString(fmt.Sprintf("%d\t%.6f\n", p.Evaluacion, p.Costo))
		}
	}
	if err := os.WriteFile(rutaDat, []byte(dat.String()), 0644); err != nil {
		return "", "", fmt.Errorf("guardando datos de convergencia: %w", err)
	}

	// --- script .gnuplot ---
	rutaDatAbs, err := filepath.Abs(rutaDat)
	if err != nil {
		return "", "", fmt.Errorf("resolviendo ruta absoluta del .dat: %w", err)
	}

	var lineasPlot []string
	indice := 0
	for _, r := range resultados {
		if len(r.Historial) == 0 {
			continue
		}
		lineasPlot = append(lineasPlot, fmt.Sprintf(
			"'%s' index %d using 1:2 with steps lw 2 title 'semilla %d'",
			filepath.ToSlash(rutaDatAbs), indice, r.Semilla,
		))
		indice++
	}

	rutaPng := filepath.Join(carpeta, fmt.Sprintf("convergencia_%s.png", marca))
	rutaPngAbs, err := filepath.Abs(rutaPng)
	if err != nil {
		return "", "", fmt.Errorf("resolviendo ruta absoluta del .png: %w", err)
	}

	var script strings.Builder
	script.WriteString("set terminal pngcairo size 1000,600 enhanced font 'Arial,10'\n")
	script.WriteString(fmt.Sprintf("set output '%s'\n", filepath.ToSlash(rutaPngAbs)))
	script.WriteString(fmt.Sprintf("set title 'Convergencia -- modo: %s'\n", modo))
	script.WriteString("set xlabel 'Numero de evaluaciones'\n")
	script.WriteString("set ylabel 'Costo (normalizado)'\n")
	script.WriteString("set grid\n")
	script.WriteString("set key outside right\n")
	script.WriteString("plot \\\n  " + strings.Join(lineasPlot, ", \\\n  ") + "\n")

	if err := os.WriteFile(rutaScript, []byte(script.String()), 0644); err != nil {
		return "", "", fmt.Errorf("guardando script de gnuplot: %w", err)
	}

	return rutaDat, rutaScript, nil
}
