package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/logan-bobo/blog-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func (apiCfg *apiConfig) readiness(w http.ResponseWriter, r *http.Request) {
	jsonPayload := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	respondWithJSON(w, 200, jsonPayload)
}

func (apiCfg *apiConfig) error(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}

func(apiCfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request) {
	var apiKey string

	authHeader := r.Header.Get("Authorization")	

	if authHeader == "" {
		respondWithError(w, 400, "No Auth header supplied")
		return
	}

	splitAuth := strings.Split(authHeader, " ")

	if len(splitAuth) > 0 {
		apiKey = splitAuth[1]
	} else {
		respondWithError(w, 400, "Incorrect Auth header format")
		return
	}

	user, err := apiCfg.DB.SelectUserAPIKey(r.Context(), apiKey)

	if err != nil {
		respondWithError(w, 400, "No user exists for API key")
		return
	}

	type getUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	response := getUserResponse{
		ID:		   user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}

	respondWithJSON(w, 200, response)
}

func (apiCfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type putUserRequest struct {
		Name string `json:"name"`
	}

	user := putUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		respondWithError(w, 400, "Invalid client request body")
		return
	}

	currentTime := time.Now()

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      user.Name,
	}

	createdUser, err := apiCfg.DB.CreateUser(r.Context(), userParams)

	if err != nil {
		respondWithError(w, 500, "Can not create user in database")
		return
	}

	type putUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}

	response := putUserResponse{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Name:      createdUser.Name,
	}

	respondWithJSON(w, 201, response)
}
