package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/gorilla/mux"
	// "github.com/gorilla/mux"
)

func GetAsistanceTypesHandler(w http.ResponseWriter, r *http.Request) {
	asistanceSql := "SELECT * FROM `asistance_types`"
	asistances := models.AsistanceTypes{}

	datos, err := db.DB.Query(asistanceSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.AsistanceType{}
		if err != nil {
			w.Write([]byte("Error en la peticion"))
		}
		datos.Scan(&dato.Asistance_type_uid, &dato.Name)
		asistances = append(asistances, dato)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(asistances)
}

func GetAsistancePlayerbyIdHandler(w http.ResponseWriter, r *http.Request) {
	assistanceSql := "SELECT asistance.player_uid, asistance.event_uid, asistance.asistance_type_uid, asistance_types.name, events.date, events.event_name FROM asistance INNER JOIN events ON asistance.event_uid = events.event_uid INNER JOIN asistance_types ON asistance.asistance_type_uid = asistance_types.asistance_type_uid WHERE asistance.player_uid ="
	vars := mux.Vars(r)

	asistance := models.Asistances{}
	datos, err := db.DB.Query(assistanceSql + "'" + vars["id"] + "'")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		log.Fatal("Error obteniendo asistencia: ", err)
		json.NewEncoder(w).Encode(err)
	}

	for datos.Next() {
		dato := models.Asistance{}
		if err != nil {
			w.Write([]byte("Error en la peticion"))
		}
		datos.Scan(&dato.Event_uid, &dato.Player_uid, &dato.Asistance_type_uid, &dato.Name, &dato.Date, &dato.Event_name)
		asistance = append(asistance, dato)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(asistance)
}
