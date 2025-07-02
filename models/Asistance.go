package models

type AsistanceType struct {
	Asistance_type_uid string `json:"asistance_type_uid"`
	Name               string `json:"asistance_type_name"`
}

type AsistanceTypes []AsistanceType
