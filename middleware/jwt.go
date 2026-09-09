package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, bool) {

	// errorVariables := godotenv.Load()
	// if errorVariables != nil {
	// 	return nil, false
	// }
	var jwtKey = []byte(os.Getenv("HASH"))

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return nil, true
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, false
	}

	return nil, true
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.UnauthorizedResponse(w, "Missing or invalid Authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, ok := ValidateJWT(tokenStr)
		if ok {
			helpers.UnauthorizedResponse(w, "Unauthorized: invalid token")
			return
		}

		// Attach user info to context
		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
