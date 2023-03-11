package main

import (
	"log"

	"github.com/retrixe/wire-transfer/core"
)

func main() {
	LoadConfig()

	// TODO: drop this
	println(core.ErrInvalidPacket)

	log.Println("Listening on port", config.Port)
}
