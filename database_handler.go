package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Afsinoz/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) UserHandler(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Email string `json:"email"`
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
	user_email := params.Email

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid,
		UpdatedAt: currentTime,
		Email:     user_email,
	})
	if err != nil {
		log.Printf("Error while creating the data base: %s", err)
		w.WriteHeader(500)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	responseWithJson(w, 201, user)

}
