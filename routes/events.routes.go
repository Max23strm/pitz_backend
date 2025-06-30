package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
)

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	db.DBconnection()
	eventsSql := "SELECT event_uid, event_type, date, event_name, event_types.type_name, events_state.event_state FROM events INNER JOIN event_types ON events.event_type = event_types.event_type_uid INNER JOIN events_state ON events.event_state_uid = events_state.event_state_uid; "
	events := models.Events{}

	datos, err := db.DB.Query(eventsSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.Event{}
		if err != nil {
			w.Write([]byte("Error en la peticion"))
		}
		datos.Scan(&dato.Event_uid, &dato.Event_type_uid, &dato.Date, &dato.Type_name, &dato.Event_name, &dato.Event_state)
		events = append(events, dato)
	}
	defer db.CerrarConexion()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func GetEventsTypesHandler(w http.ResponseWriter, r *http.Request) {
	db.DBconnection()
	eventsSql := "SELECT * FROM event_types"
	events := models.EventsTypes{}

	datos, err := db.DB.Query(eventsSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.EventType{}
		if err != nil {
			w.Write([]byte("Error en la peticion"))
		}
		datos.Scan(&dato.Event_type_uid, &dato.Type_name)
		events = append(events, dato)
	}
	defer db.CerrarConexion()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
