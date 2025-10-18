package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Coddyum/chirpy/handler"
	"github.com/Coddyum/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env fils
	godotenv.Load()

	// Create the http server
	mux := http.NewServeMux()

	// Get the db url from .env fils
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	// Open the sql db from the url
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error when openning sql db")
		return
	}

	dbQueries := database.New(db)

	apiCfg := &handler.ApiConfig{
		DB:       dbQueries,
		Platform: platform,
	}

	// this is for listen http server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", fileServer)))
	mux.Handle("/assets", fileServer)

	// Metric
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetMetricHandler)

	// Status off app health
	mux.HandleFunc("GET /api/healthz", handler.ReadlinessHandler)

	// Json
	// mux.HandleFunc("POST /api/validate_chirp", handler.ValidatedChirp)

	// Users
	mux.HandleFunc("POST /api/users", apiCfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.LoginHandler)

	mux.HandleFunc("POST /api/chirps", apiCfg.CreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerSelectAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerSelectOneChirp)

	srv.ListenAndServe()
}
