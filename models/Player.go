package models

import "database/sql"

type Player struct {
	Player_uid string `json:"player_uid"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Status     int16  `json:"status"`
	Positions  []int8 `json:"positions"`
}

type PlayerDetails struct {
	GeneralInfo Player         `json:"player_data"`
	Address     sql.NullString `json:"address"`
}
type PlayerDetailsWithAsistance struct {
	GeneralInfo Player         `json:"player_data"`
	Address     sql.NullString `json:"address"`
	Asistance   Asistances     `json:"asistance"`
}

type Positions struct {
	Position_uid string `json:"positions_uid"`
	Positions    string `json:"positions"`
}

type Players []Player

type NewPlayer struct {
	Player_uid string
	FirstName  string
	LastName   string
	Email      string
	Status     int16
	Positions  string
}
