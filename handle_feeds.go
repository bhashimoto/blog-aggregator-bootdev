package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
)

func (cfg *apiConfig) HandleFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 401, "invalid request")
		return
	}
	ctx := context.Background()
	feed, err := cfg.db.CreateFeeed(ctx, database.CreateFeeedParams{
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("unable to create feed: %s", err.Error()))
		return
	}
	respondWithJSON(w, 201, feed)
}
