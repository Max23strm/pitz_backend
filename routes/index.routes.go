package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Max23strm/pitz-backend/calendar"
	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
)

func HomeHanlder(w http.ResponseWriter, r *http.Request) {
	playersSql := "SELECT COUNT(*) AS active_players FROM players WHERE players.status = 1;"
	incomeSql := "SELECT  COALESCE(SUM(payments.amount), 0) AS monthly_income FROM payments WHERE payments.delete_flag = 0 AND payments.date BETWEEN ? AND ?;"
	expensesSql := "SELECT  COALESCE(SUM(expenses.amount), 0) AS monthly_expense FROM expenses WHERE expenses.delete_flag = 0 AND expenses.date BETWEEN ? AND ?;"

	dateStr := r.URL.Query().Get("date")

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
	currentDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	startOfMonthFormated := startOfMonth.Format("2006-01-02 15:04:05")
	endOfMonthFormated := endOfMonth.Format("2006-01-02 15:04:05")

	finalResponse := models.FinalResponse{}

	incomeRow := db.DB.QueryRow(incomeSql, startOfMonthFormated, endOfMonthFormated)
	expenseRow := db.DB.QueryRow(expensesSql, startOfMonthFormated, endOfMonthFormated)
	playersRow := db.DB.QueryRow(playersSql)

	err, ThisEvent := calendar.GetNextEvent(currentDate.Format(time.RFC3339), endOfMonth.Format(time.RFC3339))
	if err == nil {
		start := ThisEvent.Start.DateTime
		end := ThisEvent.End.DateTime
		if start == "" {
			start = ThisEvent.Start.Date
		}
		if end == "" {
			end = ThisEvent.End.Date
		}

		finalResponse.UpcomingEvent.Summary = ThisEvent.Summary
		finalResponse.UpcomingEvent.Created = ThisEvent.Created
		finalResponse.UpcomingEvent.EventType = ThisEvent.EventType
		finalResponse.UpcomingEvent.HtmlLink = ThisEvent.HtmlLink
		finalResponse.UpcomingEvent.Kind = ThisEvent.Kind
		finalResponse.UpcomingEvent.Location = ThisEvent.Location
		finalResponse.UpcomingEvent.Start = ThisEvent.Start
		finalResponse.UpcomingEvent.End = ThisEvent.End
		finalResponse.UpcomingEvent.Status = ThisEvent.Status
	} else {
		fmt.Println("---- Error de calendario ----")
		fmt.Println(err)
	}

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
	err = expenseRow.Scan(&finalResponse.Monthly_expense)
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

	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "OK",
		"data":      finalResponse,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)

}
