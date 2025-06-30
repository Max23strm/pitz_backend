package main

import (
	"net/http"

	"github.com/Max23strm/pitz-backend/routes"
	"github.com/gorilla/mux"
)

// docker run --name some-postgres -e ¨PSTGRES_USER=max -e POSTGRES_PASSWORD=secretpassword -d postgres
func main() {

	r := mux.NewRouter()

	r.HandleFunc("/loginSession", routes.LoginSession)
	r.HandleFunc("/passwordRestoration", routes.RestorePassword)

	r.HandleFunc("/home", routes.HomeHanlder)

	//USERS - PLAYERS
	r.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	r.HandleFunc("/users", routes.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")

	//EVENTS
	r.HandleFunc("/events", routes.GetEventsHandler).Methods("GET")

	//EVENT TYPES
	r.HandleFunc("/eventsTypes", routes.GetEventsTypesHandler).Methods("GET")
	// r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	// r.HandleFunc("/users", routes.DeleteUserHandler).Methods("DELETE")
	// r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")

	http.ListenAndServe(":3050", r)
}
