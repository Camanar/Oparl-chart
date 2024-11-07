package main

import (
	"log"
	"net/http"

	"github.com/Camanar/Oparl-chart/charts"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Start worker goroutines

	r := mux.NewRouter()
	r.HandleFunc("/lineChart", charts.GetAmendments).Methods("GET")

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:1090", "http://localhost:3001"}, // Add your frontend's local URL here
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          false,
	})

	// Wrap the router with the CORS handler
	handler := c.Handler(r)

	log.Println("Server starting on :1170")
	log.Fatal(http.ListenAndServe(":1170", handler))
}
