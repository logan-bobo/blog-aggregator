package main

import (
	"net/http"

	"github.com/logan-bobo/blog-aggregator/internal/auth"
	"github.com/logan-bobo/blog-aggregator/internal/database"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}

		user, err := apiCfg.DB.SelectUserAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 401, "Invalid user Key")
			return
		}

		handler(w, r, user)
	})
}
