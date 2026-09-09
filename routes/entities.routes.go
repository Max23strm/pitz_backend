package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/lib/pq"
)

func GetEntitiesByUser(w http.ResponseWriter, r *http.Request) {
	userUID, ok := r.Context().Value("userID").(string)
	if !ok || userUID == "" {
		helpers.UnauthorizedResponse(w, "Usuario no autenticado")
		return
	}

	entitiesAssigned := `SELECT
			assignment.user_uid,
			entity.entity_uid,
			entity.short_name,
			entity.name,
			entity.currency_code,
			entity.country_code,
			entity.colors
		FROM user_entities AS assignment
		INNER JOIN entities AS entity
			ON assignment.entity_uid = entity.entity_uid
		WHERE assignment.user_uid = $1
		AND entity.delete_flag = 0;`

	datos, err := db.DB.Query(entitiesAssigned, userUID)
	if err != nil {
		helpers.BadRequestResponse(w, "Error fetching entities", err, "ERROR_FETCHING")
		return
	}
	defer datos.Close()

	entities := models.EntitiesAssigned{}

	for datos.Next() {
		dato := models.EntityAssignation{}

		if err := datos.Scan(&dato.User_uid, &dato.Entity_uid, &dato.Short_name, &dato.Name, &dato.Currency_code, &dato.Country_code, pq.Array(&dato.Colors)); err != nil {
			helpers.BadRequestResponse(w, "Error fetching entities", err, "ERROR_FETCHING")
			return
		}

		entities = append(entities, dato)
	}
	if err := datos.Err(); err != nil {
		helpers.BadRequestResponse(w, "Error fetching entities", err, "ERROR_FETCHING")
		return
	}

	helpers.SuccessResponse(w, "Success", entities)
}

func InsertEntity(w http.ResponseWriter, r *http.Request) {

	if !validations.ValidateContext(w, r) {
		return
	}

	var entity = models.EntityCeate{}

	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		helpers.BadRequestResponse(w, "No information sent", err, "EMPTY_DATA")
		return
	}

	//Validacion de campos requeridos
	validationErrors := validations.EntitiesPosValidations(entity)

	if len(validationErrors) > 0 {
		message := strings.Join(validationErrors, ", ")

		helpers.BadRequestResponse(
			w,
			message,
			errors.New(message),
			"VALIDATION_ERRORS",
		)

		return
	}

	sqlString := "INSERT INTO \"entities\" ( \"short_name\", \"name\", \"currency_code\", \"country_code\", \"colors\") VALUES ($1, $2, $3, $4, $5)"

	_, err := db.DB.Exec(sqlString, entity.Short_name, entity.Name, entity.Currency_code, entity.Country_code, pq.Array(entity.Colors))
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Unable to store information",
			err,
			"STORING_ERROR",
		)
		return
	}

	helpers.CreatedResponse(
		w,
		"Stored succesfully",
		nil,
	)
}
func AssignEntity(w http.ResponseWriter, r *http.Request) {

	userUID, ok := r.Context().Value("userID").(string)
	if !ok || userUID == "" {
		helpers.UnauthorizedResponse(w, "User no log in")
		return
	}

	var entityAsignation = models.EntityAssignation{}

	if err := json.NewDecoder(r.Body).Decode(&entityAsignation); err != nil {
		helpers.BadRequestResponse(w, "No information sent", err, "EMPTY_DATA")
		return
	}

	sqlString := "INSERT INTO \"user_entities\" ( \"user_uid\", \"entity_uid\", \"assigned_by_uid\" ) VALUES ($1, $2, $3)"

	_, err := db.DB.Exec(sqlString, entityAsignation.User_uid, entityAsignation.Entity_uid, userUID)
	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Unable to store information",
			err,
			"STORING_ERROR",
		)
		return
	}

	helpers.CreatedResponse(
		w,
		"Stored succesfully",
		nil,
	)
}

func GetAllEntities(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}

	entitiesSql := "SELECT entity_uid, short_name, name, currency_code, country_code, colors, delete_flag FROM entities"
	entities := models.Entities{}

	datos, err := db.DB.Query(entitiesSql)
	if err != nil {
		helpers.ErrorResponse(w, http.StatusBadRequest, "Error obtaining entities", err)

		return
	}

	for datos.Next() {
		dato := models.Entitie{}
		err := datos.Scan(&dato.Entity_uid, &dato.Short_name, &dato.Name, &dato.Currency_code, &dato.Country_code, pq.Array(&dato.Colors), &dato.Delete_flag)
		if err != nil {
			helpers.ErrorResponse(w, http.StatusBadRequest, "Error scanning players", err)
			return
		}

		entities = append(entities, dato)
	}

	helpers.SuccessResponse(w, "", entities)
}
