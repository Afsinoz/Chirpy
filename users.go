package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Email            string    `json:"email"`
	HashedPassword   string    `json:"-"`
	ExpiresInSeconds *int      `json:"expires_in_seconds"` // * means it is optional

	Token string `json:"token"`
}
