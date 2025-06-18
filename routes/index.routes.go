package routes

import "net/http"

func HomeHanlder(w http.ResponseWriter, r *http.Request) {
	
	w.Write([]byte("Hello worldlis"))
}
