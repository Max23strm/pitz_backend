package validations

import (
	"net/http"
	"strings"

	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/middleware"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return bytes, err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateContext(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		helpers.UnauthorizedResponse(
			w,
			"Log in to continue",
		)
		return false
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		helpers.UnauthorizedResponse(
			w,
			"Log in to continue",
		)
		return false
	}

	_, err := middleware.ValidateJWT(parts[1])
	if err {
		helpers.UnauthorizedResponse(
			w,
			"Not authorized",
		)
		return false
	}

	return true
}

func GeneratePass(password string) (string, error) {
	generatedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(generatedPass), err
}
