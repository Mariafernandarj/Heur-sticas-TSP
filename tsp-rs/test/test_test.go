package test

import (
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"tsp-rs/internal/data"
	"tsp-rs/internal/model"
)

func TestCargarArchivo_Exito(t *testing.T) {
	// Crear un archivo temporal que se elimina automáticamente al terminar la prueba
	dirTemp := t.TempDir()
	rutaArchivo := filepath.Join(dirTemp, "input.tsp")

	// Escribir datos de prueba en el archivo
	contenido := "1,2,3,4,5,6,7,8,9,11,12"
	if err := os.WriteFile(rutaArchivo, []byte(contenido), 0644); err != nil {
		t.Fatalf("Error al preparar el archivo de prueba: %v", err)
	}

	// Ejecutar el método a probar
	resultado, err := data.CargarArchivo(rutaArchivo)

	// Verificaciones
	if err != nil {
		t.Fatalf("Se esperaba éxito pero ocurrió un error: %v", err)
	}

	esperado := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12}
	if !reflect.DeepEqual(resultado, esperado) {
		t.Errorf("El resultado no coincide.\nObtenido: %v\nEsperado: %v", resultado, esperado)
	}
}

func TestCargarArchivo_ArchivoNoExiste(t *testing.T) {
	_, err := data.CargarArchivo("ruta_inexistente.tsp")
	if err == nil {
		t.Error("Se esperaba un error al intentar abrir un archivo inexistente, pero err fue nil")
	}
}

func TestDistanciaNatural(t *testing.T) {
	// Ciudad 1: Tokio
	latU := 35.68500000000000227
	longU := 139.7510000000000047

	// Ciudad 2: Manila
	latV := 14.60420000000000051
	longV := 120.9819999999999994

	// Distancia esperada en metros
	distanciaEsperada := 2999396.229999999982

	// Ejecutar método
	distanciaObtenida := model.DistanciaNatural(latU, longU, latV, longV)

	// Tolerancia para errores de punto flotante
	tolerancia := 1e-7

	error := math.Abs(distanciaObtenida - distanciaEsperada)

	t.Logf("Esperada: %.18f", distanciaEsperada)
	t.Logf("Obtenida: %.18f", distanciaObtenida)
	t.Logf("Error: %.18f", error)

	// Verificar resultado
	if error > tolerancia {
		t.Errorf(
			"distancia fuera de tolerancia: esperada=%.18f, obtenida=%.18f, error=%g, tolerancia=%g",
			distanciaEsperada,
			distanciaObtenida,
			error,
			tolerancia,
		)
	}
}
