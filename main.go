package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/api"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable not set")
	}

	tokenString := os.Getenv("TOKEN_STRING")
	if platform == "" {
		log.Fatal("TOKEN_STRING env var not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	const port = "8080"
	const filepathRoot = "."

	cfg := &internal.ApiConfig{
		Database:    dbQueries,
		Platform:    platform,
		TokenString: tokenString,
	}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	apiHandler := api.NewHandler(cfg)

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(fileServerHandler))

	mux.HandleFunc("GET /admin/metrics", cfg.GetHits)
	mux.HandleFunc("POST /admin/reset", apiHandler.ResetHits)

	mux.HandleFunc("GET /api/healthz", apiHandler.HealthzHandler)

	mux.HandleFunc("POST /api/users", apiHandler.CreateNewUser)
	mux.HandleFunc("PUT /api/users", apiHandler.UpdateUserData)
	mux.HandleFunc("POST /api/login", apiHandler.Login)
	mux.HandleFunc("POST /api/refresh", apiHandler.Refresh)
	mux.HandleFunc("POST /api/revoke", apiHandler.Revoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiHandler.UpgradeUser)

	mux.HandleFunc("POST /api/chirps", apiHandler.CreateNewChirp)
	mux.HandleFunc("GET /api/chirps", apiHandler.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiHandler.GetOneChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiHandler.DeleteOneChirp)

	testServ := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server starting...\n\n")
	log.Printf("Now serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(testServ.ListenAndServe())

}
