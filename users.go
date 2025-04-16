package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Afsinoz/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Token          string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
}

func GetUser(r *http.Request, secret string) (uuid.UUID, int, string, error) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		msg := fmt.Sprintf("Invalid token error, couldn't get the token: %s", err)
		return uuid.UUID{}, 401, msg, err
	}

	reqUserID, err := auth.ValidateJWT(bearerToken, secret)
	if err != nil {
		msg := fmt.Sprintf("Unauthorized user: %s, error: %s", reqUserID, err)
		return uuid.UUID{}, 401, msg, err
	}

	return reqUserID, 0, "", nil

}
