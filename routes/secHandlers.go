package routes

import (
	"fmt"
	"net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(string)
	fmt.Fprintf(w, "Hello, authenticated user %s", userID)
}

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the public endpoint!")
}
