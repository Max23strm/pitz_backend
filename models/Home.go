package models

import "time"

type NextEvent struct {
	Event_uid  string    `json:"event_uid"`
	Event_name string    `json:"event_name"`
	Date       time.Time `json:"date"`
	Type_name  string    `json:"type_name"`
}

type FinalResponse struct {
	Monthly_income  float64    `json:"monthly_income"`
	Monthly_expense float64    `json:"monthly_expense"`
	Players_amount  float64    `json:"players_amount"`
	UpcomingEvent   *NextEvent `json:"upcoming_event"`
}
