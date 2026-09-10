package validations

import (
	"fmt"
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

func EventsPostValidations(event models.EventPost) []string {

	var validationErrors []string

	if strings.TrimSpace(event.Event_name) == "" {
		validationErrors = append(validationErrors, "El nombre es requerido.")
	}
	fmt.Println(event.Date)
	if event.Date.IsZero() {
		validationErrors = append(validationErrors, "La fecha es requerida.")
	}

	if strings.TrimSpace(event.Event_type_uid) == "" {
		validationErrors = append(validationErrors, "El tipo de evento es requerido.")
	}

	if strings.TrimSpace(event.Event_state_uid) == "" {
		validationErrors = append(validationErrors, "El estado es requerido.")
	}

	return validationErrors
}
