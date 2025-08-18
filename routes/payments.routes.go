package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/helpers"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetMonthPaymentsHandler(w http.ResponseWriter, r *http.Request) {
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

	startOfMonthFormated := startOfMonth.Format("2006-01-02 15:04:05")
	endOfMonthFormated := endOfMonth.Format("2006-01-02 15:04:05")

	paymentSQL := "SELECT payments.payment_uid, CONCAT(payer.first_name, ' ', payer.last_name) AS player_name, CONCAT(registrar.first_name, ' ', registrar.last_name) AS registered_by_name, payments.player_uid, payments.amount, payments.date, payment_type.payment_name FROM `payments` INNER JOIN players AS payer ON payments.player_uid = payer.player_uid INNER JOIN users AS registrar ON payments.registered_by_uid = registrar.user_uid 	 INNER JOIN payment_type ON payments.payment_type_uid = payment_type.payment_type_uid WHERE payments.delete_flag = 0 AND payments.date BETWEEN ? AND ?  ORDER by payments.date DESC"
	payments := models.Payments{}

	datos, err := db.DB.Query(paymentSQL, startOfMonthFormated, endOfMonthFormated)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.PaymentGeneral{}
		err := datos.Scan(&dato.Payment_uid, &dato.Player_name, &dato.Creator_name, &dato.Player_uid, &dato.Amount, &dato.Date, &dato.Payment_name)
		if err != nil {
			log.Println("Error obteniendo datos: ", err)
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}

		payments = append(payments, dato)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

func PostMonthPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	var payment models.PostPayments

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	validationErrors := validations.PaymentsPostValidations(payment)

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

	new_uuid := uuid.New()
	paymentSql := "INSERT INTO `payments` (`payment_uid`, `player_uid`, `payment_reference`, `amount`, `comment`, `date`, `payment_type_uid`, `created_at_dttm`, `updated_at_dttm`, `registered_by_uid`) VALUES (?, ?, ?, ?, ?, ?, ?, current_timestamp(), current_timestamp(), ?);"
	_, err := db.DB.Exec(paymentSql, new_uuid.String(), payment.Player_uid, payment.Payment_reference, payment.Amount, payment.Comment, payment.Date, payment.Payment_type_uid, payment.User_uid)

	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error al crear en el sql",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	uuidResponse := map[string]interface{}{
		"isSuccess":   true,
		"estado":      "Creado",
		"payment_uid": new_uuid.String(),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uuidResponse)
}

func GetPaymentTypesHandler(w http.ResponseWriter, r *http.Request) {

	paymentSQL := "SELECT * FROM payment_type"
	payment_type := models.PaymentTypes{}

	datos, err := db.DB.Query(paymentSQL)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.PaymentType{}
		err := datos.Scan(&dato.Payment_type_uid, &dato.Payment_name)
		if err != nil {
			log.Println("Error obteniendo datos: ", err)
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}

		payment_type = append(payment_type, dato)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment_type)
}

func GetPaymentByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	paymentSql := "SELECT payments.payment_uid, payments.payment_reference, payments.amount, payments.comment, payments.date, CONCAT(players.first_name, ' ', players.last_name) as player_name, payments.player_uid, payment_type.payment_name, CONCAT(creator.first_name, ' ', creator.last_name) as registered_by FROM `payments` INNER JOIN players on players.player_uid = payments.player_uid INNER JOIN users as creator on creator.user_uid = payments.registered_by_uid INNER JOIN payment_type on payment_type.payment_type_uid = payments.payment_type_uid WHERE payment_uid = ? AND payments.delete_flag = 0"

	vars := mux.Vars(r)
	paymentRow := db.DB.QueryRow(paymentSql, vars["id"])
	payment := models.PaymentById{}

	err := paymentRow.Scan(&payment.Payment_uid, &payment.Player_reference, &payment.Amount, &payment.Comment, &payment.Date, &payment.Player_name, &payment.Player_uid, &payment.Payment_name, &payment.Creator_name)
	if err != nil {
		if err == sql.ErrNoRows {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "No encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "Ok",
		"data":      payment,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}

func DeletePaymentByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	paymentSql := "UPDATE payments SET delete_flag = 1 WHERE payment_uid = ? "

	vars := mux.Vars(r)

	if vars["id"] == "undefined" || len(vars["id"]) == 0 {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Ok",
			"mensaje":   "EL id es necesario",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	paymentRow, err := db.DB.Exec(paymentSql, vars["id"])
	rowsAffected, err := paymentRow.RowsAffected()
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	if rowsAffected == 0 {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "No encontrado",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	respuesta := map[string]interface{}{
		"isSuccess":   true,
		"estado":      "Ok",
		"payment_uid": vars["id"],
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}

func GetPaymentsReport(w http.ResponseWriter, r *http.Request) {

	var paymentRequested models.PaymentFile

	if err := json.NewDecoder(r.Body).Decode(&paymentRequested); err != nil {
		respuesta := map[string]string{
			"estado":  "Error",
			"mensaje": "Error al crear",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}
	validationErrors := validations.PaymentsFileValidation(paymentRequested)

	if len(validationErrors) > 0 {

		var errors []string
		errors = append(errors, validationErrors...)

		respuesta := map[string]interface{}{
			"isSuccese": false,
			"estado":    "Error",
			"mensaje":   errors,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	// startDate, err := time.Parse("2006-01-02", paymentRequested.Start_date)

	paymentDataSql := "SELECT payments.date AS payment_date, payments.payment_uid, payer.first_name as player_name,payer.last_name as player_last_name,CONCAT(registrar.first_name, ' ', registrar.last_name) AS registered_by_name, payments.player_uid, payments.amount, payments.comment, payment_type.payment_name FROM `payments` INNER JOIN players AS payer ON payments.player_uid = payer.player_uid INNER JOIN users AS registrar ON payments.registered_by_uid = registrar.user_uid INNER JOIN payment_type ON payments.payment_type_uid = payment_type.payment_type_uid WHERE payments.delete_flag = 0 AND payments.date BETWEEN ? AND ? ORDER by payments.date ASC"
	monthlySQL := "SELECT DATE_FORMAT(date, '%m-%Y') AS month, SUM(amount) AS total FROM payments WHERE payments.delete_flag = 0 AND payments.date BETWEEN ? AND ? GROUP BY month ORDER BY month;"

	paymentsRows, err := db.DB.Query(paymentDataSql, paymentRequested.Start_date, paymentRequested.End_date)
	payments := models.PaymentFileRows{}

	if err != nil {
		if err == sql.ErrNoRows {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "No encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	for paymentsRows.Next() {
		payment := models.PaymentFileRow{}
		err := paymentsRows.Scan(&payment.Payment_date, &payment.Payment_uid, &payment.Player_name, &payment.Player_last_name, &payment.Registered_by_name, &payment.Player_uid, &payment.Amount, &payment.Comment, &payment.Payment_name)
		if err != nil {
			log.Println("Error obteniendo datos: ", err)
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}

		payments = append(payments, payment)
	}
	monthlyRows, err := db.DB.Query(monthlySQL, paymentRequested.Start_date, paymentRequested.End_date)
	monthlyPayments := models.MonthlyFileRows{}
	if err != nil {
		if err == sql.ErrNoRows {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "No encontrado",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error obteniendo pago: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}

	for monthlyRows.Next() {
		var monthStr string
		var total float64

		err := monthlyRows.Scan(&monthStr, &total)
		if err != nil {
			log.Fatal(err)
		}

		// Parse "08-2025" into time.Time
		parsedTime, err := time.Parse("01-2006", monthStr)
		if err != nil {
			log.Fatal("Failed to parse month:", err)
		}

		monthlyPayments = append(monthlyPayments, models.MonthlyFileRow{
			Month:  parsedTime,
			Amount: total,
		})

	}
	file, err := helpers.CreatePaymentExcel(payments, monthlyPayments)
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error generando archivo: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		http.Error(w, "Error writing Excel file", http.StatusInternalServerError)
		return
	}

	// Set headers to prompt a download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=report.xlsx")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(http.StatusOK)

	// Write the file to the response
	_, err = w.Write(buf.Bytes())
	if err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   "Error writing response: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)

		return
	}
}
