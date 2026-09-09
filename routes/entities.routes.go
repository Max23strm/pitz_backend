package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/gorilla/mux"
)

func GetEntitiesByUserId(w http.ResponseWriter, r *http.Request) {
	entitiesAssigned := `
		SELECT user_uid, entity_uid
		FROM user_entities
		WHERE user_uid = $1`
	vars := mux.Vars(r)

	datos, err := db.DB.Query(entitiesAssigned, vars["user_uid"])
	if err != nil {
		w.WriteHeader(http.StatusOK)
		log.Fatal(w, "Error obteniendo entidades", err)
		return
	}
	defer datos.Close()

	entities := models.EntitiesAssigned{}

	for datos.Next() {
		dato := models.EntityAssignation{}

		if err := datos.Scan(&dato.User_uid, &dato.Entity_uid); err != nil {
			w.Write([]byte("Error en la peticion"))
			return
		}

		entities = append(entities, dato)
	}
	if err := datos.Err(); err != nil {
		http.Error(w, "Error recorriendo entidades", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entities)
}
