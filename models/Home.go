package models

import (
	calendar "google.golang.org/api/calendar/v3"
)

type NextEvent struct {
	Created   string                  `json:"created"`
	EventType string                  `json:"event_type"`
	HtmlLink  string                  `json:"html_link"`
	Kind      string                  `json:"kind"`
	Location  string                  `json:"location"`
	Start     *calendar.EventDateTime `json:"start"`
	Status    string                  `json:"status"`
	Summary   string                  `json:"summary"`
	End       *calendar.EventDateTime `json:"end"`
}

type FinalResponse struct {
	Monthly_income  float64   `json:"monthly_income"`
	Monthly_expense float64   `json:"monthly_expense"`
	Players_amount  float64   `json:"players_amount"`
	UpcomingEvent   NextEvent `json:"upcoming_event"`
}
