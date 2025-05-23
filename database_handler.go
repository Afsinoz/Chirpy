package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Afsinoz/Chirpy/internal/auth"
	"github.com/Afsinoz/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) UserHandler(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// Decode the request
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error while decoding email: %s", err)
		w.WriteHeader(500)
		return
	}
	uuid := uuid.New()
	currentTime := time.Now()
	userEmail := params.Email
	HashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		msg := fmt.Sprintf("Hashing problem, %s", err)
		responseWithError(w, 500, msg, err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid,
		UpdatedAt:      currentTime,
		Email:          userEmail,
		HashedPassword: HashedPassword,
	})
	if err != nil {
		log.Printf("Error while creating the data base: %s", err)
		w.WriteHeader(500)
		return
	}

	user := User{
		ID:             dbUser.ID,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		Email:          dbUser.Email,
		HashedPassword: HashedPassword,
		IsChirpyRed:    dbUser.IsChirpyRed,
	}

	responseWithJson(w, 201, user)
}

func (cfg *apiConfig) ChirpsHandler(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	sort_query := r.URL.Query().Get("sort")

	if len(authorID) == 0 {
		dbChirps, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			msg := fmt.Sprintf("Error getting the list of chirps: %s", err)
			responseWithError(w, 500, msg, err)
			return
		}

		var responseChirps []Chirpy

		for _, chirp := range dbChirps {
			newChirpy := Chirpy{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID.UUID,
			}
			responseChirps = append(responseChirps, newChirpy)
		}
		if sort_query == "desc" {
			sort.Slice(responseChirps, func(i, j int) bool {
				return responseChirps[j].CreatedAt.Before(responseChirps[i].CreatedAt)
			})
			responseWithJson(w, 200, responseChirps)
			return
		}

		responseWithJson(w, 200, responseChirps)
		return
	}
	fmt.Println("HERE IS THE AUTHOR GETTING CHIRpy")
	authorUUID, err := uuid.Parse(authorID)
	if err != nil {
		responseWithError(w, 500, "uuid Parsing error", err)
		return
	}

	dbChirps, err := cfg.db.GetAuthorChirps(r.Context(),
		uuid.NullUUID{
			UUID:  authorUUID,
			Valid: true,
		})
	if err != nil {
		msg := fmt.Sprintf("Error getting the list of chirps: %s", err)
		responseWithError(w, 500, msg, err)
		return
	}

	var responseChirps []Chirpy

	for _, chirp := range dbChirps {
		newChirpy := Chirpy{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.UUID,
		}
		responseChirps = append(responseChirps, newChirpy)
	}

	if sort_query == "desc" {
		sort.Slice(responseChirps, func(i, j int) bool {
			return responseChirps[j].CreatedAt.Before(responseChirps[i].CreatedAt)
		})
		responseWithJson(w, 200, responseChirps)
		return
	} else {
		responseWithJson(w, 200, responseChirps)
		return
	}
}

func (cfg *apiConfig) GetChirpyByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		// Handle invalid UUID format
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
		return
	}

	dbchirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		msg := fmt.Sprintf("Chirpy not found, db error: %s", err)
		responseWithError(w, 404, msg, err)
	}

	chirpy := Chirpy{
		ID:        dbchirp.ID,
		CreatedAt: dbchirp.CreatedAt,
		UpdatedAt: dbchirp.UpdatedAt,
		Body:      dbchirp.Body,
		UserID:    dbchirp.UserID.UUID,
	}

	responseWithJson(w, 200, chirpy)
}
