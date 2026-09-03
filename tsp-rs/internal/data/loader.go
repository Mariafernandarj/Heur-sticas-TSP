package data

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Carga el archivo de entrada input-n.tsp
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

// Función para recbir archivo desde la terminal
func RecibirArchivo() error {
	rutaPtr := flag.String("path", "input.tsp", "Ruta del archivo de entrada .tsp")
	flag.Parse()

	datos, err := CargarArchivo(*rutaPtr)
	if err != nil {
		return fmt.Errorf("al cargar archivo: %w", err)
	}
	fmt.Printf("Se leyeron ls %d elementos correctamente:\n", len(datos))
	fmt.Println(datos)
	return nil
}
