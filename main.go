package main

import (
	"log"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/middleware"
	"github.com/Max23strm/pitz-backend/routes"
	"github.com/gorilla/mux"
)

const baseUrl = "/api"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Or restrict to specific domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204 is more appropriate for preflight
			return
		}

		next.ServeHTTP(w, r)
	})
}

// docker run --name some-postgres -e ¨PSTGRES_USER=max -e POSTGRES_PASSWORD=secretpassword -d postgres
func main() {

	if err := db.DBconnection(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.CerrarConexion() // ✅ Close on shutdown

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	api := r.PathPrefix(baseUrl).Subrouter()
	api.Use(middleware.AuthMiddleware)

	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.HandleFunc(baseUrl+"/loginSession", routes.LoginSession).Methods("POST")
	r.HandleFunc(baseUrl+"/passwordRestoration", routes.RestorePassword)

	r.HandleFunc(baseUrl+"/home", routes.HomeHanlder)

	//USERS - PLAYERS
	r.HandleFunc(baseUrl+"/players", routes.GetPlayersHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/players/newPlayer", routes.PostPlayerHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/players/editPlayer/{id}", routes.EditPlayerHandler).Methods("PUT")
	r.HandleFunc(baseUrl+"/players", routes.DeletePlayerHandler).Methods("DELETE")
	r.HandleFunc(baseUrl+"/players/{id}", routes.GetPlayerByIdHandler).Methods("GET")

	//EVENTS
	r.HandleFunc(baseUrl+"/events", routes.GetEventsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/events/eventById/{id}", routes.GetEventByIdHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/events/newEvent", routes.NewEventHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/events/editEvent/{id}", routes.EditEventHandler).Methods("PUT")
	r.HandleFunc(baseUrl+"/events/getEventsByMonth", routes.GetEventsByMonthHandler).Methods("GET")

	//EVENT TYPES
	r.HandleFunc(baseUrl+"/eventsTypes", routes.GetEventsTypesHandler).Methods("GET")

	//ASISTANCE TYPES
	r.HandleFunc(baseUrl+"/asistanceTypes", routes.GetAsistanceTypesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/asistanceByPlayerId/{id}", routes.GetAsistancePlayerbyIdHandler).Methods("GET")

	//PAYMENTS
	r.HandleFunc(baseUrl+"/payments", routes.GetMonthPaymentsHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/payments/new-payment", routes.PostMonthPaymentsHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/paymentsTypes", routes.GetPaymentTypesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/payments/paymentById/{id}", routes.GetPaymentByIdHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/payments/deletePayment/{id}", routes.DeletePaymentByIdHandler).Methods("DELETE")

	//Expenses
	r.HandleFunc(baseUrl+"/expenses", routes.GetMonthExpensesHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/expenses/new-expense", routes.PostNewExpenseHandler).Methods("POST")
	r.HandleFunc(baseUrl+"/expenses/{id}", routes.GetMonthExpensesByIdHandler).Methods("Get")

	//Users
	r.HandleFunc(baseUrl+"/users/", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/users/basics/{id}", routes.GetBasicUserHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc(baseUrl+"/updateUserPassword/{id}", routes.UpdateUserHandler).Methods("PUT")

	http.ListenAndServe(":3050", r)
}
