package validations

import (
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

func PlayersPostValidations(player models.PostPlayerDetails) []string {

	var validationErrors []string

	if strings.TrimSpace(player.FirstName) == "" {
		validationErrors = append(validationErrors, "El nombre es requerido.")
	}
	if strings.TrimSpace(player.LastName) == "" {
		validationErrors = append(validationErrors, "El apellido es requerido.")
	}
	if strings.TrimSpace(player.Email) == "" {
		validationErrors = append(validationErrors, "El email es requerido.")
	}
	if strings.TrimSpace(player.Sex) == "" {
		validationErrors = append(validationErrors, "El sexo es requerido.")
	}
	if player.Status < 0 || player.Status > 1 {
		validationErrors = append(validationErrors, "El estado es incorrecto.")
	}

	return validationErrors
}
