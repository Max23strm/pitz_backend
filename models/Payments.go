package models

import (
	"time"
)

type PaymentGeneral struct {
	Payment_uid  string    `json:"payment_uid"`
	Player_name  string    `json:"player_name"`
	Player_uid   string    `json:"player_uid"`
	Amount       float64   `json:"amount"`
	Date         time.Time `json:"date"`
	Payment_name string    `json:"payment_name"`
}

type Payments []PaymentGeneral

type PaymentType struct {
	Payment_type_uid string `json:"payment_type_uid"`
	Payment_name     string `json:"payment_name"`
}
type PaymentTypes []PaymentType
