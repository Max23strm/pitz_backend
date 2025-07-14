package models

import (
	"time"
)

type Player struct {
	Player_uid string `json:"player_uid"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Status     int16  `json:"status"`
	Positions  []int8 `json:"positions"`
}

type PlayerDetails struct {
	Player_uid       string    `json:"player_uid"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Status           int16     `json:"status"`
	Birth_dt         time.Time `json:"birth_dt"`
	Positions        []*int8   `json:"positions"`
	Address          *string   `json:"address"`
	Sex              string    `json:"sex"`
	BloodType        *string   `json:"blood_type"`
	Comments         *string   `json:"comments"`
	Credential       *string   `json:"credential"`
	Afiliation       *string   `json:"afiliation"`
	Curp             *string   `json:"curp"`
	Enfermedad       *string   `json:"enfermedad"`
	Phone_number     *string   `json:"phone_number"`
	Emergency_number *string   `json:"emergency_number"`
	Insurance        bool      `json:"insurance"`
	Insurance_name   *string   `json:"insurance_name"`
}
type PlayerDetailsWithAsistance struct {
	GeneralInfo Player     `json:"player_data"`
	Address     *string    `json:"address"`
	Asistance   Asistances `json:"asistance"`
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
