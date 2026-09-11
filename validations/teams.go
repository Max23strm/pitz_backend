package validations

import (
	"strings"

	"github.com/Max23strm/pitz-backend/models"
	"github.com/google/uuid"
)

func TeamsPostValidations(team models.TeamPost) []string {

	var validationErrors []string

	if strings.TrimSpace(team.Description) == "" {
		validationErrors = append(validationErrors, "Description is required")
	}

	if strings.TrimSpace(team.Entity_uid) == "" {
		validationErrors = append(validationErrors, "Entity is required")
	}

	return validationErrors
}

func TeamsAsignValidations(asignation models.PlayerAsign) []string {

	var validationErrors []string

	if strings.TrimSpace(asignation.Player_uid) == "" {
		validationErrors = append(validationErrors, "Player_uid is required")
	}
	_, err := uuid.Parse(asignation.Player_uid)
	if err != nil {
		validationErrors = append(validationErrors, "Player_uid is not correct")
	}

	if strings.TrimSpace(asignation.Team_uid) == "" {
		validationErrors = append(validationErrors, "Team_uid is required")
	}
	_, err = uuid.Parse(asignation.Team_uid)
	if err != nil {
		validationErrors = append(validationErrors, "Team_uid is not correct")
	}

	return validationErrors
}

func TeamsCatPostValidations(team models.TeamPostCategory) []string {

	var validationErrors []string

	if strings.TrimSpace(team.Description) == "" {
		validationErrors = append(validationErrors, "Description is required")
	}

	if strings.TrimSpace(team.Entity_uid) == "" {
		validationErrors = append(validationErrors, "Entity is required")
	}

	return validationErrors
}

func TeamsCatAssignValidations(team models.TeamAssignCategory) []string {

	var validationErrors []string

	if strings.TrimSpace(team.Team_uid) == "" {
		validationErrors = append(validationErrors, "Team is required")
	}

	if len(team.Categories) == 0 {
		validationErrors = append(validationErrors, "Categories are required")
	}

	for _, value := range team.Categories {
		_, err := uuid.Parse(value)
		if err != nil {
			validationErrors = append(validationErrors, "Id is not the correct format")
		}
	}

	return validationErrors
}
