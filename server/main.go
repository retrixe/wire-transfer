package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	LoadConfig()

	log.Println("Listening on port", config.Port)
	if config.HTTPS.Enabled {
		log.Fatalln(http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.HTTPS.Cert, config.HTTPS.Key, nil))
	} else {
		log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
	}
}
