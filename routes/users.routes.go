package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	db.DBconnection()
	playersSql := "SELECT players.player_uid, first_name, last_name, status, assignedpositions.positions, email FROM players  INNER JOIN assignedpositions ON players.player_uid = assignedpositions.player_uid"
	players := models.Players{}

	datos, err := db.DB.Query(playersSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.Player{}
		if err != nil {
			w.Write([]byte("Error en la peticion"))
		}
		datos.Scan(&dato.Player_uid, &dato.FirstName, &dato.LastName, &dato.Status, &dato.Positions, &dato.Email)
		players = append(players, dato)
	}
	defer db.CerrarConexion()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get specific user"))
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	// var user models.Player

	// json.NewDecoder(r.Body).Decode(&user)

	w.Write([]byte("Create user"))
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
