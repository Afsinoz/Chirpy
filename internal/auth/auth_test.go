package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name        string
		UserID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Valid case",
			UserID:      uuid.New(),
			tokenSecret: "secret key",
			expiresIn:   time.Hour,
			wantErr:     false,
		},
		{
			name:        "Zero duration",
			UserID:      uuid.New(),
			tokenSecret: "valid secret",
			expiresIn:   0,
			wantErr:     false,
		},
		{
			name:        "empty secret",
			UserID:      uuid.New(),
			tokenSecret: "",
			expiresIn:   time.Hour,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.UserID, tt.tokenSecret, tt.expiresIn)
			if tt.wantErr {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("MakeJWT failed: %v", err)
			}
			if token == "" {
				t.Fatal("Expected token to be non-empty")
			}
		})
	}

}

func TestValidateJWT(t *testing.T) {
	validSecret := "valid-secret"
	validUserID := uuid.New()

	testCases := []struct {
		name       string
		tokenSetup func() string
		secret     string
		wantErr    bool
		expectedID uuid.UUID
	}{
		{
			name: "Valid Token",
			tokenSetup: func() string {
				token, _ := MakeJWT(validUserID, validSecret, time.Hour)
				return token
			},
			secret:     validSecret,
			wantErr:    false,
			expectedID: validUserID,
		},
		{
			name: "Expired Token",
			tokenSetup: func() string {
				token, _ := MakeJWT(validUserID, validSecret, -time.Hour)
				return token
			},
			secret:     validSecret,
			wantErr:    true,
			expectedID: uuid.Nil,
		},
		{
			name: "wrong-secret",
			tokenSetup: func() string {
				token, _ := MakeJWT(validUserID, validSecret, time.Hour)
				return token
			},
			secret:     "wrong-secret",
			wantErr:    true,
			expectedID: uuid.Nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tokenString := tt.tokenSetup()

			userID, err := ValidateJWT(tokenString, tt.secret)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Expected Error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if userID != tt.expectedID {
				t.Fatalf("Expected user ID %v but got %v", tt.expectedID, userID)
			}
		})
	}

}
