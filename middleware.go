package main

import (
	"context"
	"net/http"

)

type authedHandler func(http.ResponseWriter, *http.Request, User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := getAPIKey(r)
		if err != nil {
			respondWithError(w, 401, "invalid API Key")
			return
		}
		ctx := context.Background()
		dbUser, err := cfg.db.GetUserByAPIKey(ctx, apiKey)
		if err != nil {
			respondWithError(w, 404, "user not found")
			return
		}
		user := databaseUserToUser(dbUser)
		handler(w, r, user)
	}
}

