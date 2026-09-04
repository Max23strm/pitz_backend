package helpers

import (
	"strconv"

	"github.com/Max23strm/pitz-backend/models"
)

func DefineFields(player models.PostPlayerDetails) ([]string, []interface{}) {
	fields := []string{}
	values := []interface{}{}

	idx := 1
	addField := func(col string, val interface{}) {
		fields = append(fields, col+" = $"+strconv.Itoa(idx))
		values = append(values, val)
		idx++
	}

	addField("player_uid", "")
	addField("first_name", player.FirstName)
	addField("last_name", player.LastName)

	if player.Phone_number != nil {
		addField("phone_number", *player.Phone_number)
	} else {
		addField("phone_number", nil)
	}

	if player.Emergency_number != nil {
		addField("emergency_phone", *player.Emergency_number)
	} else {
		addField("emergency_phone", nil)
	}

	addField("email", player.Email)
	addField("status", player.Status)

	// if player.Positions != nil {
	// 	values = append(values, *player.Positions)
	// } else {
	addField("positions", nil)
	// }

	addField("birth_dt", player.Birth_dt)

	if player.BloodType != nil {
		addField("blood_type", *player.BloodType)
	} else {
		addField("blood_type", nil)
	}

	if player.Comments != nil {
		addField("comments", *player.Comments)
	} else {
		addField("comments", nil)
	}

	if player.Credential != nil {
		addField("credential", *player.Credential)
	} else {
		addField("credential", nil)
	}

	if player.Address != nil {
		addField("address", *player.Address)
	} else {
		addField("address", nil)
	}

	if player.Afiliation != nil {
		addField("afiliation", *player.Afiliation)
	} else {
		addField("afiliation", nil)
	}

	addField("sex", player.Sex)

	if player.Curp != nil {
		addField("curp", *player.Curp)
	} else {
		addField("curp", nil)
	}

	if player.Enfermedad != nil {
		addField("enfermedad", *player.Enfermedad)
	} else {
		addField("enfermedad", nil)
	}

	addField("insurance", player.Insurance)

	if player.Insurance_name != nil {
		addField("insurance_name", *player.Insurance_name)
	} else {
		addField("insurance_name", nil)
	}

	return fields, values
}
