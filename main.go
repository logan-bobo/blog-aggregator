package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/logan-bobo/blog-aggregator/internal/database"
	"github.com/logan-bobo/blog-aggregator/internal/scraper"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// TODO: Move this to being an env var make const for now
const numberOfFeedsToUpdate int32 = 10

func main() {
	godotenv.Load()
	serverPort := os.Getenv("PORT")
	dbURL := os.Getenv("PG_CONN")
	fmt.Println(dbURL)

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: corsMux,
	}

	mux.HandleFunc("GET /v1/readiness", apiCfg.readiness)
	mux.HandleFunc("GET /v1/error", apiCfg.error)

	mux.HandleFunc("POST /v1/users", apiCfg.postUser)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.getUser))

	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.postFeed))
	mux.HandleFunc("GET /v1/feeds", apiCfg.getFeeds)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.postFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.deleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.getFeedFollows))

	fmt.Printf("Serving port : %v \n", serverPort)

	go scraper.Worker(numberOfFeedsToUpdate, dbQueries)

	log.Fatal(server.ListenAndServe())
}
