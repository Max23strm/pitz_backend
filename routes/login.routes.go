package routes

import "net/http"

func LoginSession(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login!"))
}

func RestorePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Should send an email!"))
}
