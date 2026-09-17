set terminal pngcairo size 1000,600 enhanced font 'Arial,10'
set output '/home/mafer/seminarioA/Heur-sticas-TSP/tsp-rs/graficas/convergencia_09-17_11-26-27.png'
set title 'Convergencia -- modo: multi'
set xlabel 'Numero de evaluaciones'
set ylabel 'Costo (normalizado)'
set grid
set key outside right
plot \
  '/home/mafer/seminarioA/Heur-sticas-TSP/tsp-rs/graficas/convergencia_09-17_11-26-27.dat' index 0 using 1:2 with steps lw 2 title 'semilla 48', \
  '/home/mafer/seminarioA/Heur-sticas-TSP/tsp-rs/graficas/convergencia_09-17_11-26-27.dat' index 1 using 1:2 with steps lw 2 title 'semilla 123459', \
  '/home/mafer/seminarioA/Heur-sticas-TSP/tsp-rs/graficas/convergencia_09-17_11-26-27.dat' index 2 using 1:2 with steps lw 2 title 'semilla 98721', \
  '/home/mafer/seminarioA/Heur-sticas-TSP/tsp-rs/graficas/convergencia_09-17_11-26-27.dat' index 3 using 1:2 with steps lw 2 title 'semilla 11'
