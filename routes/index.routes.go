package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
)

func HomeHanlder(w http.ResponseWriter, r *http.Request) {
	playersSql := "SELECT COUNT(*) AS active_players FROM players WHERE players.status = 1;"
	incomeSql := "SELECT  COALESCE(SUM(payments.amount), 0) AS monthly_income FROM payments WHERE payments.delete_flag = 0 AND payments.date BETWEEN ? AND ?;"
	eventSql := "SELECT event_uid, event_name, date, event_types.type_name FROM events INNER JOIN event_types  ON event_types.event_type_uid = events.event_type  WHERE events.date > ? ORDER BY events.date ASC LIMIT 1"
	// Get the date from the query param
	dateStr := r.URL.Query().Get("date") // e.g., "2025-06-01"

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Invalid date format. Use YYYY-MM-DD",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	currentDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	startOfMonthFormated := startOfMonth.Format("2006-01-02 15:04:05")
	endOfMonthFormated := endOfMonth.Format("2006-01-02 15:04:05")
	currentDayFormated := currentDate.Format("2006-01-02 15:04:05")

	db.DBconnection()
	finalResponse := models.FinalResponse{}

	incomeRow := db.DB.QueryRow(incomeSql, startOfMonthFormated, endOfMonthFormated)
	playersRow := db.DB.QueryRow(playersSql)
	EventRow := db.DB.QueryRow(eventSql, currentDayFormated)
	err = incomeRow.Scan(&finalResponse.Monthly_income)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo ingreso mensual: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	err = playersRow.Scan(&finalResponse.Players_amount)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo jugadores activos: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	finalResponse.UpcomingEvent = &models.NextEvent{}
	err = EventRow.Scan(&finalResponse.UpcomingEvent.Event_uid, &finalResponse.UpcomingEvent.Event_name, &finalResponse.UpcomingEvent.Date, &finalResponse.UpcomingEvent.Type_name)
	if err != nil {
		if err == sql.ErrNoRows {
			finalResponse.UpcomingEvent = nil
		} else {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "Error obteniendo proximo evento: " + err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(respuesta)

			return
		}
	}
	defer db.CerrarConexion()
	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "OK",
		"data":      finalResponse,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)

}
