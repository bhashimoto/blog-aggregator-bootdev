package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleFeedsGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	feeds, err := cfg.db.GetAllFeeds(ctx)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("could not retrieve feeds: %s", err.Error()))
		return
	}

	respondWithJSON(w, 200, feeds)
}

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
	feed, err := cfg.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
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

	ff, err := cfg.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	ret := struct {
		Feed database.Feed		`json:"feed"`
		FeedFollow database.FeedFollow	`json:"feed_follow"`
	}{
		Feed: feed,
		FeedFollow: ff,
	}

	respondWithJSON(w, 201, ret)
}
