package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/gorilla/mux"
)

func GetPlayersHandler(w http.ResponseWriter, r *http.Request) {
	db.DBconnection()
	playersSql := "SELECT players.player_uid, first_name, last_name, status, assignedpositions.positions, email FROM players  INNER JOIN assignedpositions ON players.player_uid = assignedpositions.player_uid"
	players := models.Players{}

	datos, err := db.DB.Query(playersSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.Player{}
		var positionsStr string
		err := datos.Scan(&dato.Player_uid, &dato.FirstName, &dato.LastName, &dato.Status, &positionsStr, &dato.Email)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			log.Fatal("Error obteniendo datos: ", err)
			json.NewEncoder(w).Encode(err)
		}
		err = json.Unmarshal([]byte(positionsStr), &dato.Positions)

		if err != nil {
			w.WriteHeader(http.StatusOK)
			log.Fatal("Transformando posiciones: ", err)
			json.NewEncoder(w).Encode(err)
		}

		players = append(players, dato)
	}
	defer db.CerrarConexion()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func GetPlayerByIdHandler(w http.ResponseWriter, r *http.Request) {
	playerSql := "SELECT players.player_uid, players.first_name, players.last_name, players.email, players.status, assignedpositions.positions, players.address FROM players INNER JOIN assignedpositions ON players.player_uid = assignedpositions.player_uid WHERE players.player_uid = ?"
	db.DBconnection()
	vars := mux.Vars(r)

	playerRow := db.DB.QueryRow(playerSql, vars["id"])

	player := models.PlayerDetails{}

	err := playerRow.Scan(&player.GeneralInfo.Player_uid, &player.GeneralInfo.FirstName, &player.GeneralInfo.LastName, &player.GeneralInfo.Email, &player.GeneralInfo.Status, &player.GeneralInfo.Positions, &player.Address)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		log.Fatal("Error obteniendo jugador: ", err)
		json.NewEncoder(w).Encode(err)
	}

	defer db.CerrarConexion()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func PostPlayerHandler(w http.ResponseWriter, r *http.Request) {
	// var user models.Player

	// json.NewDecoder(r.Body).Decode(&user)

	w.Write([]byte("Create user"))
}

func DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
