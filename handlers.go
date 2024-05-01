package main

import (
	"encoding/json"
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

func (apiCfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type getUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}

	response := getUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}

	respondWithJSON(w, 200, response)
}

func (apiCfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type postUserRequest struct {
		Name string `json:"name"`
	}

	user := postUserRequest{}
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

func (apiCfg *apiConfig) postFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type postFeedRequest struct {
		Name string `json:"name"`
		URL  string `json:"URL"`
	}

	feedData := postFeedRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&feedData)

	if err != nil {
		respondWithError(w, 400, "malformed request")
	}

	currentTime := time.Now()

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      feedData.Name,
		Url:       feedData.URL,
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), feedParams)

	if err != nil {
		respondWithError(w, 500, "Can not create feed in database")
		return
	}

	type postFeedResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}

	response := postFeedResponse{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}

	respondWithJSON(w, 201, response)
}

func (apiCfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.SelectAllFeeds(r.Context())

	if err != nil {
		respondWithError(w, 500, "Can not select feeds")
	}

	type returnFeed struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAT time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}

	returnFeeds := []returnFeed{}

	for _, feed := range feeds {
		returnFeeds = append(returnFeeds, returnFeed{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAT: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		})
	}

	respondWithJSON(w, 200, returnFeeds)
}

func (apiCfg *apiConfig) postFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type postFeedFollow struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	feedFollowData := postFeedFollow{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&feedFollowData)

	if err != nil {
		respondWithError(w, 400, "Malformed request")
		return
	}

	currentTime := time.Now()

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feedFollowData.FeedId,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), feedFollowParams)

	if err != nil {
		respondWithError(w, 500, "Can not create feed follow in database")
	}

	type returnFeedFollow struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	response := returnFeedFollow{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	}

	respondWithJSON(w, 201, response)
}
