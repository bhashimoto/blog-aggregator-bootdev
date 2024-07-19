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


func (cfg *apiConfig) HandleFeedFollowsGet(w http.ResponseWriter, r *http.Request, user User) {
	ctx := context.Background()
	ffs, err := cfg.db.GetFeedFollowsFromUser(ctx, user.ID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("unable to retrieve data: %s", err.Error()))
		return
	}

	respondWithJSON(w, 200, ffs)
}

func (cfg *apiConfig) HandleFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user User) {
	feedFollowIDString := r.PathValue("feedFollowID")

	if feedFollowIDString == "" {
		respondWithError(w, 401, "invalid url")
		return
	}
	ffID, err := uuid.Parse(feedFollowIDString)

	if err != nil {
		respondWithError(w, 401, "invalid uuid")
		return
	}

	ctx := context.Background()
	feedFollow, err:= cfg.db.GetFeedFollowByID(ctx, ffID)
	if err != nil {
		respondWithError(w, 404, "feed follow not found")
		return
	}

	if user.ID != feedFollow.UserID {
		respondWithError(w, http.StatusForbidden, "not allowed to delete other user feed")
		return
	}

	err = cfg.db.DeleteFeedFollow(ctx, ffID)
	if err != nil {
		respondWithError(w, 500, "unable to remove feed follow")
		return
	}

	respondWithJSON(w, 204, "")


}


func (cfg *apiConfig) HandleFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 401, "invalid request")
		return	
	}

	ctx := context.Background()
	feedFollowDB, err := cfg.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		FeedID: params.FeedID,
		UserID: user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, 500, "unable to create feed follow")
		return
	}
	feedFollow := databaseFeedFollowToFeedFollow(feedFollowDB)
	respondWithJSON(w, 201, feedFollow)
}
