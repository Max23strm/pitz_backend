package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPlayersHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	playersSql := "SELECT players.player_uid, first_name, last_name, status, email FROM players"
	players := models.Players{}

	datos, err := db.DB.Query(playersSql)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo jugadores: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	for datos.Next() {
		dato := models.Player{}
		err := datos.Scan(&dato.Player_uid, &dato.FirstName, &dato.LastName, &dato.Status, &dato.Email)
		if err != nil {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "Error obteniendo datos: " + err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(respuesta)

			return
		}

		players = append(players, dato)
	}

	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "OK",
		"data":      players,
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(respuesta)
}

func GetPlayerByIdHandler(w http.ResponseWriter, r *http.Request) {
	playerSql := "SELECT players.player_uid, players.first_name, players.last_name, players.email, players.status, players.address, players.birth_dt, players.comments, players.blood_type, players.afiliation, players.sex, players.curp, players.enfermedad, players.phone_number, players.emergency_phone, players.insurance, players.insurance_name FROM players WHERE players.player_uid = ?"

	vars := mux.Vars(r)

	playerRow := db.DB.QueryRow(playerSql, vars["id"])

	player := models.PlayerDetails{}

	err := playerRow.Scan(&player.Player_uid, &player.FirstName, &player.LastName, &player.Email, &player.Status, &player.Address, &player.Birth_dt, &player.Comments, &player.BloodType, &player.Afiliation, &player.Sex, &player.Curp, &player.Enfermedad, &player.Phone_number, &player.Emergency_number, &player.Insurance, &player.Insurance_name)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo jugador: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uuidResponse)
}

func EditPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var event models.PutPlayerDetails
	player_uid := vars["id"]
	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		fmt.Println(err)
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al obtener datos",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	fields := []string{}
	values := []interface{}{}

	if event.FirstName != nil {
		fields = append(fields, "first_name = ?")
		values = append(values, *event.FirstName)
	}
	if event.LastName != nil {
		fields = append(fields, "last_name = ?")
		values = append(values, *event.LastName)
	}
	if event.Email != nil {
		fields = append(fields, "email = ?")
		values = append(values, *event.Email)
	}
	if event.Status != nil {
		fields = append(fields, "status = ?")
		values = append(values, *event.Status)
	}
	if event.Birth_dt != nil {
		fields = append(fields, "birth_dt = ?")
		values = append(values, *event.Birth_dt)
	}
	if event.Address != nil {
		fields = append(fields, "address = ?")
		values = append(values, *event.Address)
	}
	if event.Sex != nil {
		fields = append(fields, "sex = ?")
		values = append(values, *event.Sex)
	}
	if event.BloodType != nil {
		fields = append(fields, "blood_type = ?")
		values = append(values, *event.BloodType)
	}
	if event.Comments != nil {
		fields = append(fields, "comments = ?")
		values = append(values, *event.Comments)
	}
	if event.Credential != nil {
		fields = append(fields, "credential = ?")
		values = append(values, *event.Credential)
	}
	if event.Afiliation != nil {
		fields = append(fields, "afiliation = ?")
		values = append(values, *event.Afiliation)
	}
	if event.Curp != nil {
		fields = append(fields, "curp = ?")
		values = append(values, *event.Curp)
	}
	if event.Enfermedad != nil {
		fields = append(fields, "enfermedad = ?")
		values = append(values, *event.Enfermedad)
	}
	if event.Phone_number != nil {
		fields = append(fields, "phone_number = ?")
		values = append(values, *event.Phone_number)
	}
	if event.Emergency_number != nil {
		fields = append(fields, "emergency_phone = ?")
		values = append(values, *event.Emergency_number)
	}
	if event.Insurance != nil {
		fields = append(fields, "insurance = ?")
		values = append(values, *event.Insurance)
	}
	if event.Insurance_name != nil {
		fields = append(fields, "insurance_name = ?")
		values = append(values, *event.Insurance_name)
	}

	if len(fields) == 0 {
		respuesta := map[string]interface{}{
			"isSuccess":  false,
			"estado":     "No fields to update",
			"player_uid": player_uid,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	values = append(values, player_uid)

	query := fmt.Sprintf("UPDATE players SET %s WHERE player_uid = ?", strings.Join(fields, ", "))

	result, err := db.DB.Exec(query, values...)

	if err != nil {
		errRes := map[string]interface{}{
			"isSuccess":  false,
			"estado":     err,
			"player_uid": player_uid,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errRes)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Jugador no encontrado",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	uuidResponse := map[string]interface{}{
		"isSuccess":  true,
		"estado":     "Editado",
		"player_uid": vars["id"],
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uuidResponse)
}

func DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
