package validations

import (
	"encoding/json"
	"net/http"
	"strings"

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
		w.WriteHeader(http.StatusUnauthorized)
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Inicie sesión para continuar",
		}
		json.NewEncoder(w).Encode(respuesta)
		return false
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Inicie sesión para continuar",
		}
		json.NewEncoder(w).Encode(respuesta)
		return false
	}

	_, err := middleware.ValidateJWT(parts[1])
	if err {
		w.WriteHeader(http.StatusUnauthorized)
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Sin autorizacón",
		}
		json.NewEncoder(w).Encode(respuesta)
		return false
	}

	return true
}

func GeneratePass(password string) (string, error) {
	generatedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(generatedPass), err
}
