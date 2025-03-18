package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hashedBytes)
	return hashedPassword, nil
}

func CheckPasswordHash(password string, hash string) error {
	passwordBytes := []byte(password)
	hashBytes := []byte(hash)

	if err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes); err != nil {
		return err
	}
	return nil
}
