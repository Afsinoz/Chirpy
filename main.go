package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Afsinoz/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var apiCfg apiConfig

	// Database Connection
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	apiCfg.polkaKey = os.Getenv("POLKA_KEY")

	apiCfg.secret = os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	apiCfg.platform = "dev"
	apiCfg.db = dbQueries

	const port = "8080"
	const filePathRoot = "./templates"
	mux := http.NewServeMux()
	apiCfg.secret = os.Getenv("SECRET")
	// Handler
	strpr := http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(strpr))
	// Registering handlers

	mux.HandleFunc("GET /api/healthz", ReadinessHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.RequestHandler)

	mux.HandleFunc("GET /api/chirps", apiCfg.ChirpsHandler)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpyByIDHandler)

	mux.HandleFunc("POST /admin/reset", apiCfg.ResetNumberRequestHandler)

	mux.HandleFunc("POST /api/users", apiCfg.UserHandler)

	mux.HandleFunc("POST /api/chirps", apiCfg.ChirpCreateHandler)

	mux.HandleFunc("POST /api/login", apiCfg.LoginHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshHandler)

	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeTokenHandler)

	mux.HandleFunc("PUT /api/users", apiCfg.UpdateUserHandler)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirpHandler)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.PolkaHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Error of ListenAndServer: %s", err)
	}
}
