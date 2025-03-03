package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/Afsinoz/Chirpy/internal/database"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("contentType", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK, ready!"))

}

// API CONFIG
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) RequestHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform == "dev" {
		w.Header().Set("contentType", "text/html")
		w.WriteHeader(http.StatusOK)
		s := fmt.Sprintf("<html> <body> <h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
		w.Write([]byte(s))
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}

func (cfg *apiConfig) ResetNumberRequestHandler(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Error while truncating users table: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("contentType", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("All resetted!"))
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling Json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func responseWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Err string `json:"error"`
	}

	responseWithJson(w, code, errorResponse{Err: msg})

}

func ChirpyValidationHandler(w http.ResponseWriter, r *http.Request) {

	type chirpyText struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirpyText{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, 400, "Chirpy is too long!", nil)
		return
	}

	s1 := "kerfuffle"
	s2 := "sharbert"
	s3 := "fornax"

	newText := []string{}

	splittedText := strings.Split(params.Body, " ")
	var wordCheck string
	for _, word := range splittedText {
		wordCheck = strings.ToLower(word)
		if wordCheck == s1 || wordCheck == s2 || wordCheck == s3 {
			newText = append(newText, "****")
		} else {
			newText = append(newText, word)
		}
	}
	newBody := strings.Join(newText, " ")

	clndbdy := returnValue{
		CleanedBody: newBody,
	}
	responseWithJson(w, http.StatusOK, clndbdy)

}
