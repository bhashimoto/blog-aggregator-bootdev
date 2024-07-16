package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	parameters := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		return
	}

	newUser := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: parameters.Name,
	}
	user, err := cfg.db.CreateUser(ctx, newUser)
	if err != nil {
		respondWithError(w, 500, "error creating user")
		return
	}


	respondWithJSON(w, http.StatusCreated, user)
}
