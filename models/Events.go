package models

type Event struct {
	Event_uid      string `json:"event_uid"`
	Event_type_uid string `json:"event_type_uid"`
	Date           string `json:"event_date"`
	Event_name     string `json:"event_name"`
	Type_name      string `json:"event_type_name"`
	Event_state    string `json:"event_state"`
}

type EventType struct {
	Event_type_uid string `json:"event_type_uid"`
	Type_name      string `json:"evebt_type_name"`
}

type Events []Event
type EventsTypes []EventType
