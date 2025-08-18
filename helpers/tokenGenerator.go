package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, time.Time, error) {
	// errorVariables := godotenv.Load()
	// if errorVariables != nil {
	// 	panic(errorVariables)
	// }
	var jwtKey = []byte(os.Getenv("HASH"))

	expirationTime := time.Now().Add(48 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}
