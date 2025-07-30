package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
)

func getUserFromDB(username string) (*models.UserFromDb, error) {
	var user models.UserFromDb
	db.DBconnection()

	query := `SELECT user_uid, username, hashed_password FROM users WHERE username = ? Or email= ? LIMIT 1;`
	row := db.DB.QueryRow(query, username, username)

	err := row.Scan(&user.User_uid, &user.User, &user.HashedPassword)
	if err != nil {
		return nil, err // could be sql.ErrNoRows
	}

	defer db.CerrarConexion()
	return &user, nil
}

func LoginSession(w http.ResponseWriter, r *http.Request) {

	var creds models.LoginCred
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error al obtener datos" + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	user, err := getUserFromDB(creds.User)

	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo usuario",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	if !validations.CheckPassword(creds.Password, strings.TrimSpace(user.HashedPassword)) {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Contraseña inválida",
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	token, expiration, err := helpers.GenerateJWT(user.User_uid)

	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error interno",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	tokenResponse := map[string]interface{}{
		"isSuccess":  true,
		"estado":     "Ok",
		"token":      token,
		"expiration": expiration,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)

}

func RestorePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Should send an email!"))
}
