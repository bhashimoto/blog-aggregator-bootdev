package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Environment variables
	godotenv.Load()
	
	// Set up database
	dbURL := os.Getenv("DB_CONN")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	// Set up config
	cfg := apiConfig{
		db: dbQueries,
	}

	// Set up server
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}


	log.Println("calling FetchFeeds")
	go cfg.FetchFeedsRoutine(10, 60*time.Second)
	log.Println("called FetchFeeds")


	mux.HandleFunc("GET /v1/healthz", HandleHealthz)
	mux.HandleFunc("GET /v1/err", HandleError)	
	
	mux.HandleFunc("POST /v1/users", cfg.HandleUserCreate)	
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.HandleUsersGet))	
	
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.HandleFeedsCreate))	
	mux.HandleFunc("GET /v1/feeds", cfg.HandleFeedsGet)
	mux.HandleFunc("DELETE /v1/feeds", cfg.middlewareAuth(cfg.HandleFeedsDelete))

	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.HandleFeedFollowsGet))	
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.HandleFeedFollowsCreate))	
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.HandleFeedFollowsDelete))


	log.Println("Starting server at port", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}


