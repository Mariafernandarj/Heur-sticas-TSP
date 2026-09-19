# Heurísticas para el Problema del Agente Viajero (TSP)

Proyecto de implementación y experimentación de heurísticas para resolver el **Problema del Agente Viajero (Traveling Salesman Problem, TSP)**.

## Requisitos

Para ejecutar el proyecto es necesario tener instalado:

* [Go](https://go.dev/)
* [GNUplot](http://www.gnuplot.info/)
* Git

## 1. Clonar el repositorio

Clona el repositorio utilizando:

```bash
git clone https://github.com/Mariafernandarj/Heur-sticas-TSP
```

Después, entra a la carpeta del proyecto:

```bash
cd Heur-sticas-TSP/tsp-rs
```

## 2. Compilar el proyecto

Desde la carpeta `tsp-rs`, ejecuta:

```bash
go build ./cmd/tsp
```

Si la compilación termina correctamente, el proyecto estará listo para ejecutarse.

## 3. Ejecutar el programa

### Una sola semilla

Para ejecutar el algoritmo utilizando una sola semilla, utiliza el parámetro `-semilla`.

Por ejemplo:

```bash
go run ./cmd/tsp -path ./data/input-40.tsp -semilla 1
```

También puedes utilizar otro archivo de entrada:

```bash
go run ./cmd/tsp -path ./data/input-150.tsp -semilla 2
```

### Múltiples semillas

Si deseas ejecutar el algoritmo con varias semillas, utiliza la opción `-multi` y especifica las semillas mediante `-semillas`.

Por ejemplo:

```bash
go run ./cmd/tsp -path ./data/input-40.tsp -multi -semillas 4,5,6
```

O utilizando el archivo de 150 ciudades:

```bash
go run ./cmd/tsp -path ./data/input-150.tsp -multi -semillas 1,2,3
```

## 4. Resultados

Después de ejecutar el programa se genera un reporte en formato **JSON**.

### Una sola semilla

Los resultados se almacenan en:

```text
resultados.json
```

### Múltiples semillas

Cuando se ejecuta el programa con varias semillas, los resultados se almacenan en:

```text
resultados_multi.json
```

Estos archivos contienen la información obtenida durante las ejecuciones del algoritmo.

Los resultados se almacenan tambien en archivos **TXT** .

Se pueden consultar en este formato en la carpeta:
```text
resultados
```


## 5. Generación de gráficas

El proyecto genera archivos `.gnuplot` que pueden utilizarse para visualizar los resultados.

Para generar una gráfica, utiliza:

```bash
gnuplot graficas/nombre_archivo.gnuplot
```

Por ejemplo:

```bash
gnuplot graficas/convergencia_09-19_00-15-25.gnuplot
```

El nombre del archivo `.gnuplot` dependerá de la ejecución y de las gráficas generadas.
