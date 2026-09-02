package main

import (
       "log"
       "tsp-rs/internal/data"
       )

func main() {
     db, err := data.InicioDB("tsp.sql")
     if err != nil {
     	log.Fatal("Error iniciando la BD:", err)
     }
     defer db.Close()

     ciudades, err := data.GetCiudades(db)
     if err != nil {
     	log.Fatal("Error consultando ciudades", err)
     }

     for _, c := range ciudades {
     	 log.Println(c.Nombre, c.Poblacion)
     }
}
