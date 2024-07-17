package main

import (
	"context"
	"net/http"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := getAPIKey(r)
		if err != nil {
			respondWithError(w, 401, "invalid API Key")
			return
		}
		ctx := context.Background()
		user, err := cfg.db.GetUserByAPIKey(ctx, apiKey)
		if err != nil {
			respondWithError(w, 404, "user not found")
			return
		}
		handler(w, r, user)
	}
}

