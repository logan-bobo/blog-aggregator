package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (apiCfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type putUserRequest struct {
		Name string `json:"name"`
	}

	type putUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}

	user := putUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		// Should this be a 500, it could be an issue with the request body?
		// Can I handle this better?
		// Should this log on serverside
		respondWithError(w, 500, "Can not parse request body")
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
		fmt.Println(err)
		respondWithError(w, 500, "Can not create user in database")
		return
	}

	respondWithJSON(w, 201, createdUser)
}
