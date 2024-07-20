package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
)

func (cfg *apiConfig) HandlePostsGet(w http.ResponseWriter, r *http.Request, user User) {
	limitString := r.URL.Query().Get("limit")
	limit := 10
	if limitString != "" {
		num, err := strconv.Atoi(limitString)
		if err == nil {
			limit = num
		}
	}
	posts, err := cfg.db.GetPostsByUser(context.Background(), 
					database.GetPostsByUserParams{
						ID: user.ID,
						Limit: int32(limit),
					})
	if err != nil {
		respondWithError(w, 500, "unable to retrieve posts")
		return
	}
	respondWithJSON(w, 200, posts)
}
