package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Afsinoz/Chirpy/internal/database"
	"github.com/google/uuid"
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

func chirpyValidate(w http.ResponseWriter, chirp string) string {

	if len(chirp) > 140 {
		responseWithError(w, 400, "Chirpy is too long!", nil)
		return ""
	}

	s1 := "kerfuffle"
	s2 := "sharbert"
	s3 := "fornax"

	newText := []string{}

	splittedText := strings.Split(chirp, " ")
	var wordCheck string
	for _, word := range splittedText {
		wordCheck = strings.ToLower(word)
		if wordCheck == s1 || wordCheck == s2 || wordCheck == s3 {
			newText = append(newText, "****")
		} else {
			newText = append(newText, word)
		}
	}

	cleanChirp := strings.Join(newText, " ")
	return cleanChirp

}

func (cfg *apiConfig) ChirpCreateHandler(w http.ResponseWriter, r *http.Request) {

	type reqChirpy struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := reqChirpy{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	newBody := chirpyValidate(w, params.Body)

	uuid_chirp := uuid.New()
	currentTime := time.Now()



	dbChirpy, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid_chirp,
		UpdatedAt: currentTime,
		Body:      newBody,
		UserID:    uuid.NullUUID { 
			UUID: params.UserID,
			Valid: true,
		},
	})
	if err != nil {
		log.Printf("Error while creating chirpy: %s", err)
		w.WriteHeader(500)
		return
	}

	chirpy := Chirpy{
		ID:        dbChirpy.ID,
		CreatedAt: dbChirpy.CreatedAt,
		UpdatedAt: dbChirpy.UpdatedAt,
		Body:      dbChirpy.Body,
		UserID:    dbChirpy.UserID.UUID,
	}

	responseWithJson(w, 201, chirpy)

}
