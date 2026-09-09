package models

type Entities struct {
	Entity_uid    string `json:"entity_uid"`
	Name          string `json:"name"`
	Short_name    string `json:"short_name"`
	Country_code  string `json:"country_code"`
	Currency_code string `json:"currency_code"`
	Entity_name   string `json:"entity_name"`
}

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
	User_uid   string `json:"user_uid"`
	Entity_uid string `json:"entity_uid"`
}

type EntitiesAssigned []EntityAssignation

type EntityUnassign struct {
	User_entity_uid string `json:"user_entity_uid"`
}
