package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	
	dbURL := os.Getenv("DB_CONN")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	cfg := apiConfig{
		db: dbQueries,
	}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}


	mux.HandleFunc("GET /v1/healthz", HandleHealthz)
	mux.HandleFunc("GET /v1/err", HandleError)	
	mux.HandleFunc("POST /v1/users", cfg.HandleUserCreate)	

	log.Println("Starting server at port", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

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
