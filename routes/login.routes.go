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

	query := `SELECT user_uid, username, hashed_password FROM users WHERE username = $1 Or email= $2 LIMIT 1;`
	row := db.DB.QueryRow(query, username, username)

	err := row.Scan(&user.User_uid, &user.User, &user.HashedPassword)
	if err != nil {
		return nil, err // could be sql.ErrNoRows
	}

	return &user, nil
}

func LoginSession(w http.ResponseWriter, r *http.Request) {

	var creds models.LoginCred
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		helpers.BadRequestResponse(
			w,
			"Error obtainig data",
			err,
			"NO_DATA",
		)
		return
	}

	user, err := getUserFromDB(creds.User)

	if err != nil {
		helpers.InternalServerErrorResponse(
			w,
			"Error validating credentials",
		)
		return
	}

	if !validations.CheckPassword(creds.Password, strings.TrimSpace(user.HashedPassword)) {
		helpers.ForbiddenResponse(
			w,
			"Error validating credentials",
		)
		return
	}
	token, expiration, err := helpers.GenerateJWT(user.User_uid)

	if err != nil {
		helpers.InternalServerErrorResponse(
			w,
			"Error validating credentials",
		)
		return
	}

	tokenResponse := map[string]interface{}{
		"token":      token,
		"expiration": expiration,
	}
	helpers.SuccessResponse(
		w,
		"loged in",
		tokenResponse,
	)

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(tokenResponse)

}

func RestorePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Should send an email!"))
}
