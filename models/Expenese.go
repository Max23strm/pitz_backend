package models

import "time"

type PostExpenses struct {
	Reason        string    `json:"reason"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Registered_by string    `json:"registered_by"`
	Assigned_uid  string    `json:"assigned_uid"`
}

type MonthlyExpenses struct {
	Assigned_to string    `json:"assigned_to"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Reason      string    `json:"reason"`
	Expense_uid string    `json:"expense_uid"`
}

type MonthlyExpensesGroup []MonthlyExpenses
