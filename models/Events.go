package models

type Event struct {
	Event_uid      string `json:"event_uid"`
	Event_type_uid string `json:"event_type_uid"`
	Date           string `json:"event_date"`
	Event_name     string `json:"event_name"`
	Type_name      string `json:"event_type_name"`
	Event_state    string `json:"event_state"`
}

type EventDetail struct {
	Event_uid      string `json:"event_uid"`
	Event_type_uid string `json:"event_type_uid"`
	Date           string `json:"event_date"`
	Event_name     string `json:"event_name"`
	Type_name      string `json:"event_type_name"`
	Event_state    int8   `json:"event_state"`
	Address        string `json:"address"`
	Coordinates    string `json:"coordinates"`
}

type EventType struct {
	Event_type_uid string `json:"event_type_uid"`
	Type_name      string `json:"evebt_type_name"`
}

type Asistance struct {
	Event_uid          string `json:"event_uid"`
	Player_uid         string `json:"player_uid"`
	Asistance_type_uid string `json:"asistance_type_uid"`
	Name               string `json:"asistance_name"`
	Date               string `json:"event_date"`
	Event_name         string `json:"event_name"`
}

type Asistances []Asistance

type Events []Event
type EventsTypes []EventType
