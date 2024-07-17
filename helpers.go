package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func getAPIKey(r *http.Request) (string, error) {
	for name, values := range r.Header {
		for _, value := range values {
			log.Println(name, value)
		}
	}
	fullAuth := r.Header.Get("Authorization")
	if fullAuth == "" {
		return "", errors.New("missing Authorization header")
	}
	keyParts := strings.Split(fullAuth, " ")
	if len(keyParts) != 2 {
		return "", errors.New("invalid Authorization header format")
	}
	if keyParts[0] != "ApiKey" {
		return "", errors.New("invalid key")
	}
	apiKey := keyParts[1]
	return apiKey, nil

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marhsaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	err := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	respondWithJSON(w, code, err)
}

func HandleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJSON(w, 200, resp)
}

func HandleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}
