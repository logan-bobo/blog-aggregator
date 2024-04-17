package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	serverPort := os.Getenv("PORT")

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: corsMux,
	}

	mux.HandleFunc("GET /v1/readiness", apiReadiness)
	mux.HandleFunc("GET /v1/error", apiError)

	fmt.Printf("Serving port : %v \n", serverPort)

	log.Fatal(server.ListenAndServe())
}
