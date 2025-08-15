package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Max23strm/pitz-backend/db"
	"github.com/Max23strm/pitz-backend/models"
	"github.com/Max23strm/pitz-backend/validations"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func PostNewExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	var expense models.PostExpenses

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		respuesta := map[string]interface{}{
			"isSuccess": false,
			"estado":    "Error",
			"mensaje":   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	validationErrors := validations.ExpenesePostValidations(expense)

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
	expenseSql := "INSERT INTO expenses(expense_uid, assigned_uid, reason, amount, created_at_dttm, updated_at_dttm, registered_by_uid, date) VALUES (?, ?, ?, ?, current_timestamp(), current_timestamp(), ?, ?);"
	_, err := db.DB.Exec(expenseSql, new_uuid.String(), expense.Assigned_uid, expense.Reason, expense.Amount, expense.Registered_by, expense.Date)

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

func GetMonthExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	dateStr := r.URL.Query().Get("date")

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

	expensesSQL := "SELECT expenses.expense_uid, expenses.reason, CONCAT(assigned.first_name, ' ', assigned.last_name) AS assigned_to,  expenses.amount,  expenses.date FROM expenses INNER JOIN users AS assigned ON expenses.assigned_uid = assigned.user_uid INNER JOIN users AS registrar ON expenses.registered_by_uid = registrar.user_uid  WHERE expenses.delete_flag = 0 AND expenses.date BETWEEN ? AND ?  ORDER by expenses.date DESC"
	expenses := models.MonthlyExpensesGroup{}

	datos, err := db.DB.Query(expensesSQL, startOfMonthFormated, endOfMonthFormated)
	if err != nil {
		w.Write([]byte("Error en la peticion"))
	}

	for datos.Next() {
		dato := models.MonthlyExpenses{}
		err := datos.Scan(&dato.Expense_uid, &dato.Reason, &dato.Assigned_to, &dato.Amount, &dato.Date)
		if err != nil {
			respuesta := map[string]interface{}{
				"isSuccess": false,
				"estado":    "Error",
				"mensaje":   "Error obteniendo datos",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(respuesta)
			return
		}

		expenses = append(expenses, dato)
	}

	respuesta := map[string]interface{}{
		"isSuccess": true,
		"estado":    "Ok",
		"data":      expenses,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}
func GetMonthExpensesByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	expenseSql := "SELECT expenses.assigned_uid as assigned_to_uid, CONCAT(assigned.first_name, ' ', assigned.last_name) as assigned_to, expenses.reason, expenses.amount, expenses.registered_by_uid, expenses.date FROM expenses INNER JOIN users on users.user_uid = expenses.registered_by_uid INNER JOIN users as assigned on assigned.user_uid = expenses.assigned_uid WHERE expense_uid = ?"
	vars := mux.Vars(r)

	expenseRow := db.DB.QueryRow(expenseSql, vars["id"]) // e.g., "2025-06-01"

	expense := models.GetExpenseId{}

	err := expenseRow.Scan(&expense.Assigned_to_uid, &expense.Assigned_to, &expense.Reason, &expense.Amount, &expense.Registered_by_uid, &expense.Date)
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
		"data":      expense,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}
func DeleteExpenseByIdHandler(w http.ResponseWriter, r *http.Request) {
	if !validations.ValidateContext(w, r) {
		return
	}
	expenseSql := "UPDATE expenses SET delete_flag = 1 WHERE expense_uid = ? "

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

	expenseRow, err := db.DB.Exec(expenseSql, vars["id"])
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
	rowsAffected, err := expenseRow.RowsAffected()
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
