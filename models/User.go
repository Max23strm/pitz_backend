package models

type PasswordUser struct {
	New_password string `json:"new_password"`
}

type BasicUser struct {
	User_uid string `json:"user_uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
type FullUser struct {
	User_uid   string `json:"user_uid"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
