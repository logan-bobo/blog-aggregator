package main

import (
	"net/http"
)

func apiReadiness(w http.ResponseWriter, r *http.Request) {
	jsonPayload := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	respondWithJSON(w, 200, jsonPayload)
}

func apiError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}
