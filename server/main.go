package main

import (
	"log"
	"net/http"
	"strconv"
)

// TODO: Implement precaching.
// TODO: Implement upload resume after it is added to the spec.

func main() {
	LoadConfig()

	// Register endpoints.
	http.HandleFunc("/upload", uploadEndpoint)
	// TODO: Implement GET /info
	// TODO: Implement GET /download/direct
	// TODO: Implement GET /download/proxied

	// TODO: Implement the UDP proxying protocol.

	log.Println("Listening on port", config.Port)
	if config.HTTPS.Enabled {
		log.Fatalln(http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.HTTPS.Cert, config.HTTPS.Key, nil))
	} else {
		log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
	}
}
