package models

type LoginCred struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type UserFromDb struct {
	User_uid       string `json:"user_uid"`
	User           string `json:"user"`
	HashedPassword string `json:"password"`
}
