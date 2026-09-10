package models

type Entitie struct {
	Entity_uid    string   `json:"entity_uid"`
	Name          string   `json:"name"`
	Short_name    string   `json:"short_name"`
	Country_code  string   `json:"country_code"`
	Currency_code string   `json:"currency_code"`
	Colors        []string `json:"colors"`
	Logo          string   `json:"logo,omitempty"`
	Delete_flag   int      `json:"delete_flag"`
}

type Entities []Entitie

type EntityDetail struct {
	Entity_uid     string     `json:"entity_uid"`
	Name           string     `json:"name"`
	Short_name     string     `json:"short_name"`
	Entity_name    string     `json:"entity_name"`
	Country_code   string     `json:"country_code"`
	Currency_code  string     `json:"currency_code"`
	Entity_colors  []string   `json:"entitiy_colors"`
	Teams_included []Teams    `json:"teams_included"`
	Users_assigned []AllUsers `json:"users_assigned"`
}

type EntityCeate struct {
	Name          string   `json:"name"`
	Short_name    string   `json:"short_name"`
	Country_code  string   `json:"country_code"`
	Currency_code string   `json:"currency_code"`
	Colors        []string `json:"colors"`
	Logo          string   `json:"logo,omitempty"`
}
type EntityPost struct {
	Name           string     `json:"name"`
	Short_name     string     `json:"short_name"`
	Entity_name    string     `json:"entity_name"`
	Country_code   string     `json:"country_code"`
	Currency_code  string     `json:"currency_code"`
	Teams_included []Teams    `json:"teams_included"`
	Users_assigned []AllUsers `json:"users_assigned"`
}

type EntityUpdate struct {
	Name           *string     `json:"name,omitempty"`
	Short_name     *string     `json:"short_name,omitempty"`
	Entity_name    *string     `json:"entity_name,omitempty"`
	Country_code   *string     `json:"country_code,omitempty"`
	Currency_code  *string     `json:"currency_code,omitempty"`
	Teams_included *[]Teams    `json:"teams_included,omitempty"`
	Users_assigned *[]AllUsers `json:"users_assigned,omitempty"`
}

type EntityAssignation struct {
	User_uid      string   `json:"user_uid"`
	Entity_uid    string   `json:"entity_uid"`
	Name          *string  `json:"name,omitempty"`
	Short_name    *string  `json:"short_name,omitempty"`
	Country_code  *string  `json:"country_code,omitempty"`
	Currency_code *string  `json:"currency_code,omitempty"`
	Colors        []string `json:"colors"`
	Logo          string   `json:"logo,omitempty"`
}

type EntitiesAssigned []EntityAssignation

type EntityUnassign struct {
	User_entity_uid string `json:"user_entity_uid"`
}
