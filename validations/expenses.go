package validations

import (
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

func ExpenesePostValidations(expene models.PostExpenses) []string {
	var validationErrors []string

	if strings.TrimSpace(expene.Reason) == "" {
		validationErrors = append(validationErrors, "El nombre de el jugador es requerido.")
	}

	if strings.TrimSpace(expene.Assigned_uid) == "" {
		validationErrors = append(validationErrors, "El tipo de pago es requerido.")
	}

	if expene.Date.IsZero() {
		validationErrors = append(validationErrors, "La fecha es requerida.")
	}

	if expene.Amount == 0 || expene.Amount < 0 {
		validationErrors = append(validationErrors, "El monto debe ser mayor a 0.")
	}

	return validationErrors
}
