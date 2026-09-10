package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Max23strm/pitz-backend/calendar"
	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	eventsSql := "SELECT event_uid, event_type, date, event_name, event_types.type_name, events_state.event_state FROM events INNER JOIN event_types ON events.event_type = event_types.event_type_uid INNER JOIN events_state ON events.event_state_uid = events_state.event_state_uid; "
	events := models.Events{}

	datos, err := db.DB.Query(eventsSql)
	if err != nil {
		helpers.InternalServerErrorResponse(w, "Petition error")
	}

	for datos.Next() {
		dato := models.Event{}
		err := datos.Scan(&dato.Event_uid, &dato.Event_type_uid, &dato.Date, &dato.Event_name, &dato.Type_name, &dato.Event_state)
		if err != nil {
			helpers.BadRequestResponse(w, "Error scanning", err, "DATA_ERROR")
			return
		}
		events = append(events, dato)
	}

	helpers.SuccessResponse(w, "succes", events)
}

func GetEventByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	eventsSql := "SELECT event_uid, event_type, date, event_types.type_name, event_name, events_state.event_state, address, coordinates FROM events INNER JOIN event_types ON events.event_type = event_types.event_type_uid INNER JOIN events_state ON events.event_state_uid = events_state.event_state_uid WHERE event_uid = $1"
	vars := mux.Vars(r)

	eventData := db.DB.QueryRow(eventsSql, vars["id"])

	currentEvent := models.EventDetail{}

	eventData.Scan(&currentEvent.Event_uid, &currentEvent.Event_type_uid, &currentEvent.Date, &currentEvent.Type_name, &currentEvent.Event_name, &currentEvent.Event_state, &currentEvent.Address, &currentEvent.Coordinates)

	helpers.SuccessResponse(w, "succes", currentEvent)
}

func GetEventsTypesHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	eventsSql := "SELECT * FROM event_types"
	events := models.EventsTypes{}

	datos, err := db.DB.Query(eventsSql)
	if err != nil {
		helpers.InternalServerErrorResponse(w, "Error requesting")
	}

	for datos.Next() {
		dato := models.EventType{}
		err := datos.Scan(&dato.Event_type_uid, &dato.Type_name)
		if err != nil {
			helpers.BadRequestResponse(w, "Error requesting", err, "ERROR_OBTAINING")
			return
		}
		events = append(events, dato)
	}

	helpers.SuccessResponse(w, "success", events)
}

func NewEventHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	var event models.EventPost
	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		helpers.BadRequestResponse(w, "Error obtaining data", err, "ERROR_DATA")
		return
	}
	validationErrors := validations.EventsPostValidations(event)

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

	new_uuid := uuid.New()
	eventsSql := "INSERT INTO \"events\"( event_uid, \"event_type\", \"date\", \"event_name\", \"event_state_uid\", \"address\", \"coordinates\") VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := db.DB.Exec(eventsSql, new_uuid.String(), event.Event_type_uid, event.Date, event.Event_name, event.Event_state_uid, event.Address, event.Coordinates)

	if err != nil {
		helpers.BadRequestResponse(
			w,
			"Error creating",
			err,
			"STORING_ERROR",
		)
		return
	}

	uuidResponse := map[string]string{
		"event_uuid": new_uuid.String(),
	}
	helpers.CreatedResponse(w, "success", uuidResponse)
}

func EditEventHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	vars := mux.Vars(r)
	event_uid := vars["id"]

	var event models.EventUpdate

	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		helpers.BadRequestResponse(w, "Error obtainig data", err, "ERROR_DATA")
		return
	}

	fields := []string{}
	values := []interface{}{}

	idx := 1
	addField := func(col string, val interface{}) {
		fields = append(fields, fmt.Sprintf("%s = $%d", col, idx))
		values = append(values, val)
		idx++
	}

	if event.Event_name != nil {
		addField("event_name", *event.Event_name)
	}
	if event.Event_type_uid != nil {
		addField("event_type", *event.Event_type_uid)
	}
	if event.Date != nil {
		addField("date", *event.Date)
	}
	if event.Event_state_uid != nil {
		addField("event_state_uid", *event.Event_state_uid)
	}
	if event.Address != nil {
		addField("address", *event.Address)
	}
	if event.Coordinates != nil {
		addField("coordinates", *event.Coordinates)
	}

	if len(fields) == 0 {
		helpers.BadRequestResponse(w, "No changes detected", errors.New(""), "ERROR_EDITING")
		return
	}

	values = append(values, event_uid)

	// Final query
	query := fmt.Sprintf("UPDATE events SET %s WHERE event_uid = $%d", strings.Join(fields, ", "), idx)

	result, err := db.DB.Exec(query, values...)

	if err != nil {
		helpers.BadRequestResponse(w, "No changes detected", err, "ERROR_EDITING")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		helpers.BadRequestResponse(w, "No event", err, "ERROR_EDITING")
		return
	}

	helpers.SuccessResponse(w, "success", event_uid)
}

func GetEventsByMonthHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	dateStr := r.URL.Query().Get("date") // e.g., "2025-06-01"

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		helpers.BadRequestResponse(w, "Invalid date format. Use YYYY-MM-DD", err, "INVALID_DATE")
		return
	}

	// Get the first day of the month
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Get the last day of the month by going to the first day of the next month and subtracting a day
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	var googleEvents []map[string]interface{}

	startOfMonthFormated := startOfMonth.Format("2006-01-02T15:04:05Z")
	endOfMonthFormated := endOfMonth.Format("2006-01-02T15:04:05Z")

	err, fetchedEvents := calendar.GetEventsByMonth(startOfMonthFormated, endOfMonthFormated)
	if err != nil {
		helpers.InternalServerErrorResponse(w, "Error getting events")
		return
	}

	if len(fetchedEvents) == 0 {
		helpers.SuccessResponse(w, "success", googleEvents)
		return
	}
	for _, item := range fetchedEvents {
		start := item.Start.DateTime
		end := item.End.DateTime
		if start == "" {
			start = item.Start.Date
		}
		if end == "" {
			end = item.End.Date
		}

		dato := map[string]interface{}{
			"google_id":  item.Id,
			"kind":       item.Kind,
			"summary":    item.Summary,
			"location":   item.Location,
			"event_type": item.EventType,
			"start":      start,
			"end":        end,
			"link":       item.HtmlLink,
		}

		googleEvents = append(googleEvents, dato)
	}

	helpers.SuccessResponse(w, "success", googleEvents)
}

func GetEventsStatesHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	statesSql := "SELECT * FROM events_state"
	eventsStates := models.EventStates{}

	datos, err := db.DB.Query(statesSql)
	if err != nil {
		helpers.InternalServerErrorResponse(w, "Error requesting")
	}

	for datos.Next() {
		dato := models.EventState{}
		err := datos.Scan(&dato.Event_state_uid, &dato.Event_state)
		if err != nil {
			helpers.BadRequestResponse(w, "Error requesting", err, "ERROR_OBTAINING")
			return
		}
		eventsStates = append(eventsStates, dato)
	}

	helpers.SuccessResponse(w, "success", eventsStates)
}
