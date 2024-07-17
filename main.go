package main

import (
	"database/sql"
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
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.HandleUsersGet))	


	log.Println("Starting server at port", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}


