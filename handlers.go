package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("contentType", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK, ready!"))

}

// API CONFIG
type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) RequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("contentType", "text/html")
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf("<html> <body> <h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(s))
}

func (cfg *apiConfig) ResetNumberRequestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("contentType", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("All resetted!"))
}

func ChirpyValidationHandler(w http.ResponseWriter, r *http.Request) {

	type chirpyText struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirpyText{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		w.WriteHeader(500)
		return
	}

	type ResponseBody struct {
		Valid bool   `json:"valid"`
		Err   string `json:"error"`
	}

	if len(params.Body) <= 140 {
		vld := ResponseBody{
			Valid: true,
			Err:   "",
		}

		dat, err := json.Marshal(vld)
		if err != nil {
			log.Printf("Error Marshalling Json: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)

	} else {
		vld := ResponseBody{
			Valid: false,
			Err:   "Chirp is too long",
		}
		dat, err := json.Marshal(vld)
		if err != nil {
			log.Printf("Error Marshalling Json: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
	}

}
