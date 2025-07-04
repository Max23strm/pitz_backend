package main

import (
	"net/http"

	"github.com/Max23strm/pitz-backend/routes"
	"github.com/gorilla/mux"
)

const baseUrl = "/api"

// docker run --name some-postgres -e ¨PSTGRES_USER=max -e POSTGRES_PASSWORD=secretpassword -d postgres
func main() {

	r := mux.NewRouter()

	r.HandleFunc(baseUrl+"/loginSession", routes.LoginSession)
	r.HandleFunc(baseUrl+"/passwordRestoration", routes.RestorePassword)

	r.HandleFunc(baseUrl+"/home", routes.HomeHanlder)

	//USERS - PLAYERS
	r.HandleFunc(baseUrl+"/players", routes.GetPlayersHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/players", routes.PostPlayerHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/players", routes.DeletePlayerHandler).Methods("DELETE")
	r.HandleFunc(baseUrl+"/players/{id}", routes.GetPlayerByIdHandler).Methods("GET")

	//EVENTS
	r.HandleFunc(baseUrl+"/events", routes.GetEventsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/events/{id}", routes.GetEventByIdHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/events/newEvent", routes.NewEventHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/events/editEvent/{id}", routes.EditEventHandler).Methods("PUT")

	//EVENT TYPES
	r.HandleFunc(baseUrl+"/eventsTypes", routes.GetEventsTypesHandler).Methods("GET")

	//ASISTANCE TYPES
	r.HandleFunc(baseUrl+"/asistanceTypes", routes.GetAsistanceTypesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/asistanceByPlayerId/{id}", routes.GetAsistancePlayerbyIdHandler).Methods("GET")
	// r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	// r.HandleFunc("/users", routes.DeleteUserHandler).Methods("DELETE")
	// r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")

	http.ListenAndServe(":3050", r)
}
