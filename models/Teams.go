package models

type Teams struct {
	Team_uid      string `json:"team_uid"`
	Description   string `json:"description"`
	Entity_uid    string `json:"entity_uid"`
	Country_code  string `json:"country_code"`
	Currency_code string `json:"currency_code"`
	Entity_name   string `json:"entity_name"`
}

type TeamsDetails struct {
	Team_uid         string   `json:"team_uid"`
	Description      string   `json:"description"`
	Entity_uid       string   `json:"entity_uid"`
	Country_code     string   `json:"country_code"`
	Currency_code    string   `json:"currency_code"`
	Entity_name      string   `json:"entity_name"`
	Players_assigned []Player `json:"players_assigned"`
}
