package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) NumberRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("contentType", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(s))
}

func (cfg *apiConfig) ResetNumberRequestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("contentType", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("All resetted!"))
}

func main() {
	const port = "8080"
	var apiCfg apiConfig
	mux := http.NewServeMux()
	// Handler
	fs := http.FileServer(http.Dir("."))

	strpr := http.StripPrefix("/app/", fs)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(strpr))
	// Registering handlers

	mux.HandleFunc("GET /healthz", ReadinessHandler)

	mux.HandleFunc("GET /metrics", apiCfg.NumberRequestHandler)
	mux.HandleFunc("POST /reset", apiCfg.ResetNumberRequestHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Errorf("Error of ListenAndServer", err)
	}

}
