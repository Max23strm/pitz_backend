package models

import (
	"time"
)

type PaymentGeneral struct {
	Payment_uid  string    `json:"payment_uid"`
	Player_name  string    `json:"player_name"`
	Creator_name string    `json:"creator_name"`
	Player_uid   string    `json:"player_uid"`
	Amount       float64   `json:"amount"`
	Date         time.Time `json:"date"`
	Payment_name string    `json:"payment_name"`
}

type PostPayments struct {
	Player_uid        string    `json:"player_uid"`
	Payment_reference *string   `json:"payment_reference"`
	Amount            float64   `json:"amount"`
	Comment           *string   `json:"comment"`
	Date              time.Time `json:"date"`
	Payment_type_uid  string    `json:"payment_type_uid"`
	User_uid          string    `json:"user_uid"`
}

type Payments []PaymentGeneral

type PaymentType struct {
	Payment_type_uid string `json:"payment_type_uid"`
	Payment_name     string `json:"payment_name"`
}
type PaymentTypes []PaymentType

type PaymentById struct {
	Payment_uid      string    `json:"payment_uid"`
	Player_reference *string   `json:"payment_reference"`
	Amount           string    `json:"amount"`
	Comment          *string   `json:"comment"`
	Date             time.Time `json:"date"`
	Player_name      string    `json:"player_name"`
	Player_uid       string    `json:"player_uid"`
	Payment_name     string    `json:"payment_name"`
	Creator_name     string    `json:"registered_by"`
}
type PaymentUid struct {
	Payment_uid string `json:"payment_uid"`
}
