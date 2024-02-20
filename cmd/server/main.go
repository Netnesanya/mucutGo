package main

import (
	"github.com/gorilla/handlers"
	"log"
	"mucutGo/internal/api"
	"net/http"
)

func main() {
	router := api.NewRouter()
	const address = ":8080"

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*", "https://netnesanya.github.io"}) // Adjust the port to match your Tauri app's port
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	// Wrap your router with CORS middleware
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))

	log.Printf("Listening on %s", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal(err)
	}
}
