package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "./templates"
	var apiCfg apiConfig
	mux := http.NewServeMux()
	// Handler

	strpr := http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(strpr))
	// Registering handlers

	mux.HandleFunc("GET /api/healthz", ReadinessHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.RequestHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetNumberRequestHandler)

	mux.HandleFunc("POST /api/validate_chirp", ChirpyValidationHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Errorf("Error of ListenAndServer", err)
	}

}
