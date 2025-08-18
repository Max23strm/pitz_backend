package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Max23strm/pitz-backend/calendar"
	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	eventsSql := "SELECT event_uid, event_type, date, event_name, event_types.type_name, events_state.event_state FROM events INNER JOIN event_types ON events.event_type = event_types.event_type_uid INNER JOIN events_state ON events.event_state_uid = events_state.event_state_uid; "
	events := models.Events{}

	datos, err := db.DB.Query(eventsSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.Event{}
		err := datos.Scan(&dato.Event_uid, &dato.Event_type_uid, &dato.Date, &dato.Event_name, &dato.Type_name, &dato.Event_state)
		if err != nil {
			respuesta := map[string]string{
				"estado":  "Error",
				"mensaje": "Error obteniendo datos",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		events = append(events, dato)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func GetEventByIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Que hace aca?")
	eventsSql := "SELECT event_uid, event_type, date, event_types.type_name, event_name, events_state.event_state, address, coordinates FROM events INNER JOIN event_types ON events.event_type = event_types.event_type_uid INNER JOIN events_state ON events.event_state_uid = events_state.event_state_uid WHERE event_uid = ?"
	vars := mux.Vars(r)

	eventData := db.DB.QueryRow(eventsSql, vars["id"])

	currentEvent := models.EventDetail{}

	eventData.Scan(&currentEvent.Event_uid, &currentEvent.Event_type_uid, &currentEvent.Date, &currentEvent.Type_name, &currentEvent.Event_name, &currentEvent.Event_state, &currentEvent.Address, &currentEvent.Coordinates)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentEvent)
}

func GetEventsTypesHandler(w http.ResponseWriter, r *http.Request) {
	eventsSql := "SELECT * FROM event_types"
	events := models.EventsTypes{}

	datos, err := db.DB.Query(eventsSql)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.EventType{}
		err := datos.Scan(&dato.Event_type_uid, &dato.Type_name)
		if err != nil {
			respuesta := map[string]string{
				"estado":  "Error",
				"mensaje": "Error obteniendo datos",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		events = append(events, dato)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func NewEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.EventPost
	//FALTA IMPLEMENTAR EL RESTO
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al crear",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	validationErrors := validations.EventsPostValidations(event)

	if len(validationErrors) > 0 {

		var errors []string
		errors = append(errors, validationErrors...)

		respuesta := map[string]interface{}{
			"estado":  "Error",
			"mensaje": errors,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	new_uuid := uuid.New()
	eventsSql := "INSERT INTO `events`( event_uid, `event_type`, `date`, `event_name`, `event_state_uid`, `address`, `coordinates`) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := db.DB.Exec(eventsSql, new_uuid.String(), event.Event_type_uid, event.Date, event.Event_name, event.Event_state_uid, event.Address, event.Coordinates)

	if err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al crear",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	uuidResponse := map[string]string{
		"event_uuid": new_uuid.String(),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uuidResponse)
}

func EditEventHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	event_uid := vars["id"]

	var event models.EventUpdate

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

	if event.Event_name != nil {
		fields = append(fields, "event_name = ?")
		values = append(values, *event.Event_name)
	}
	if event.Event_type_uid != nil {
		fields = append(fields, "event_type = ?")
		values = append(values, *event.Event_type_uid)
	}
	if event.Date != nil {
		fields = append(fields, "date = ?")
		values = append(values, *event.Date)
	}
	if event.Event_state_uid != nil {
		fields = append(fields, "event_state_uid = ?")
		values = append(values, *event.Event_state_uid)
	}
	if event.Address != nil {
		fields = append(fields, "address = ?")
		values = append(values, *event.Address)
	}
	if event.Coordinates != nil {
		fields = append(fields, "coordinates = ?")
		values = append(values, *event.Coordinates)
	}

	if len(fields) == 0 {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "No hay campos para actualizar",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	values = append(values, event_uid)

	// Final query
	query := fmt.Sprintf("UPDATE events SET %s WHERE event_uid = ?", strings.Join(fields, ", "))

	result, err := db.DB.Exec(query, values...)

	if err != nil {
		fmt.Println(err)
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al editar",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Evento no encontrado",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	uuidResponse := map[string]string{
		"event_uuid": event_uid,
		"estado":     "Editado con éxito",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uuidResponse)
}

func GetEventsByMonthHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date") // e.g., "2025-06-01"

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
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
		eventsResponse := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   err.Error(),
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(eventsResponse)
	}

	if len(fetchedEvents) == 0 {
		eventsResponse := map[string]interface{}{
			"isSuccess": true,
			"events":    googleEvents,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(eventsResponse)
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

	eventsResponse := map[string]interface{}{
		"isSuccess": true,
		"events":    googleEvents,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(eventsResponse)
}
