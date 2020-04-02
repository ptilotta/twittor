package main

import (
	"log"

	"github.com/ptilotta/twittor/bd"
	"github.com/ptilotta/twittor/handlers"
)

func main() {
	if bd.ChequeoConnection() == 0 {
		log.Fatal("Sin conexión a la BD")
		return
	}
	handlers.Manejadores()
}
