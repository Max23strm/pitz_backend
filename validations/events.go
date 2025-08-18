package validations

import (
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

func EventsPostValidations(event models.EventPost) []string {

	var validationErrors []string

	if strings.TrimSpace(event.Event_name) == "" {
		validationErrors = append(validationErrors, "El nombre es requerido.")
	}

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

func EventsFileValidation(event models.EventFile) []string {

	var validationErrors []string

	if strings.TrimSpace(event.File_type) == "" {
		validationErrors = append(validationErrors, "El tipo de archivo es requerido.")
	}
	if event.File_type != "excel" && event.File_type != "pdf" {
		validationErrors = append(validationErrors, "El tipo de archivo es incorrecto: "+event.File_type)
	}

	if event.Start_date.IsZero() {
		validationErrors = append(validationErrors, "La fecha de incio es requerida.")
	}
	if event.End_date.IsZero() {
		validationErrors = append(validationErrors, "La fecha final es requerida.")
	}

	if event.Start_date.After(event.End_date) {
		validationErrors = append(validationErrors, "La fecha final debe ser despues de la inicial")
	}

	return validationErrors
}
