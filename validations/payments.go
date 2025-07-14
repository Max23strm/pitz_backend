package validations

import (
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

func PaymentsPostValidations(payment models.PostPayments) []string {

	var validationErrors []string

	if strings.TrimSpace(payment.Player_uid) == "" {
		validationErrors = append(validationErrors, "El nombre de el jugador es requerido.")
	}

	if strings.TrimSpace(payment.Payment_type_uid) == "" {
		validationErrors = append(validationErrors, "El tipo de pago es requerido.")
	}

	if payment.Date.IsZero() {
		validationErrors = append(validationErrors, "La fecha es requerida.")
	}

	if payment.Amount == 0 || payment.Amount < 0 {
		validationErrors = append(validationErrors, "El monto debe ser mayor a 0.")
	}

	return validationErrors
}
