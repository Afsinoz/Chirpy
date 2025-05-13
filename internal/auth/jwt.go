package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO, uuid.UUID should change with uuid.NullUUID.UUID, bc of the uuid package
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", nil
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	}
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Parse with Claims error: %s, expiresAt %v", err, claims.ExpiresAt)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid Token")
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Getting UUID of the subject error: %s", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}
	return userUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	parts := strings.Fields(authHeader)
	if len(parts) < 2 {
		return "", fmt.Errorf("Invalid token")
	}
	tokenString := parts[1]
	return tokenString, nil
}
