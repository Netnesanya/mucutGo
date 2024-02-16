package main

import (
	"log"
	"mucutGo/internal/api"
	"net/http"
)

func main() {
	router := api.NewRouter()
	const address = ":8080"

	log.Printf("Listening on %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal(err)
	}
}
