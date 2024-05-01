package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/logan-bobo/blog-aggregator/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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


	fmt.Printf("Serving port : %v \n", serverPort)

	log.Fatal(server.ListenAndServe())
}
