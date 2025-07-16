package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPlayersHandler(w http.ResponseWriter, r *http.Request) {
	db.DBconnection()
	playersSql := "SELECT players.player_uid, first_name, last_name, status, email FROM players"
	players := models.Players{}

	datos, err := db.DB.Query(playersSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.Player{}
		err := datos.Scan(&dato.Player_uid, &dato.FirstName, &dato.LastName, &dato.Status, &dato.Email)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			log.Fatal("Error obteniendo datos: ", err)
			json.NewEncoder(w).Encode(err)
		}

		players = append(players, dato)
	}
	defer db.CerrarConexion()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func GetPlayerByIdHandler(w http.ResponseWriter, r *http.Request) {
	playerSql := "SELECT players.player_uid, players.first_name, players.last_name, players.email, players.status, players.address, players.birth_dt, players.comments, players.blood_type, players.afiliation, players.sex, players.curp, players.enfermedad, players.phone_number, players.emergency_phone, players.insurance, players.insurance_name FROM players INNER JOIN assignedpositions ON players.player_uid = assignedpositions.player_uid WHERE players.player_uid = ?"
	db.DBconnection()
	vars := mux.Vars(r)

	playerRow := db.DB.QueryRow(playerSql, vars["id"])

	player := models.PlayerDetails{}

	err := playerRow.Scan(&player.Player_uid, &player.FirstName, &player.LastName, &player.Email, &player.Status, &player.Address, &player.Birth_dt, &player.Comments, &player.BloodType, &player.Afiliation, &player.Sex, &player.Curp, &player.Enfermedad, &player.Phone_number, &player.Emergency_number, &player.Insurance, &player.Insurance_name)
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
	var player models.PostPlayerDetails
	new_uuid := uuid.New()

	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		fmt.Println(err)
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al obtener datos",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	//Validacion de campos requeridos
	validationErrors := validations.PlayersPostValidations(player)

	if len(validationErrors) > 0 {

		var errors []string
		errors = append(errors, validationErrors...)

		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   errors,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	sqlString := "INSERT INTO `players` (`player_uid`, `first_name`, `last_name`, `phone_number`, `emergency_phone`, `email`, `status`, `positions`, `birth_dt`, `blood_type`, `comments`, `credential`, `address`, `afiliation`, `sex`, `curp`, `enfermedad`, `insurance`, `insurance_name`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	db.DBconnection()
	_, err := db.DB.Exec(sqlString, new_uuid.String(), player.FirstName, player.LastName, player.Phone_number, player.Emergency_number, player.Email, player.Status, nil, player.Birth_dt, player.BloodType, player.Comments, player.Credential, player.Address, player.Afiliation, player.Sex, player.Curp, player.Enfermedad, player.Insurance, player.Insurance_name)

	if err != nil {
		fmt.Println(err)
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error al enviar informaciín",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	uuidResponse := map[string]interface{}{
		"isSuccess":  true,
		"estado":     "Creado",
		"player_uid": new_uuid.String(),
	}
	defer db.CerrarConexion()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uuidResponse)
}

func DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
