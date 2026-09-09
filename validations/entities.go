package validations

import (
	"regexp"
	"strings"

	"github.com/Max23strm/pitz-backend/models"
)

var HexColorRegex = regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

func EntitiesPosValidations(entity models.EntityCeate) []string {
	var validationErrors []string

	if strings.TrimSpace(entity.Name) == "" {
		validationErrors = append(validationErrors, "Name is required")
	}
	if strings.TrimSpace(entity.Short_name) == "" {
		validationErrors = append(validationErrors, "Short name is required")
	}
	if len(entity.Colors) == 0 {
		validationErrors = append(validationErrors, "No colors where assigned")
	}
	if len(entity.Colors) > 3 {
		validationErrors = append(validationErrors, "Too many colors where assigned")
	}
	for _, col := range entity.Colors {
		result := HexColorRegex.MatchString(col)
		if !result {
			validationErrors = append(validationErrors, col+" is not a valid color")
		}
	}

	return validationErrors
}
