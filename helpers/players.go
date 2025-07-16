package helpers

import (
	"strconv"

	"github.com/Max23strm/pitz-backend/models"
)

func DefineFields(player models.PostPlayerDetails) ([]string, []interface{}) {
	fields := []string{}
	values := []interface{}{}

	fields = append(fields, "player_uid = ?")

	fields = append(fields, "firstName = ?")
	values = append(values, player.FirstName)

	fields = append(fields, "last_name = ?")
	values = append(values, player.LastName)

	fields = append(fields, "phone_number = ?")
	if player.Phone_number != nil {
		values = append(values, *player.Phone_number)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "emergency_phone = ?")
	if player.Emergency_number != nil {
		values = append(values, *player.Emergency_number)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "email = ?")
	values = append(values, player.Email)

	fields = append(fields, "status = ?")
	values = append(values, player.Status)

	fields = append(fields, "positions = ?")
	// if player.Positions != nil {
	// 	values = append(values, *player.Positions)
	// } else {
	values = append(values, nil)
	// }

	fields = append(fields, "birth_dt = ?")
	values = append(values, player.Birth_dt)

	fields = append(fields, "blood_type = ?")
	if player.BloodType != nil {
		values = append(values, *player.BloodType)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "comments = ?")
	if player.Comments != nil {
		values = append(values, *player.Comments)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "credential = ?")
	if player.Credential != nil {
		values = append(values, *player.Credential)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "address = ?")
	if player.Address != nil {
		values = append(values, *player.Address)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "afiliation = ?")
	if player.Afiliation != nil {
		values = append(values, *player.Afiliation)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "sex = ?")
	values = append(values, *&player.Sex)

	fields = append(fields, "curp = ?")
	if player.Curp != nil {
		values = append(values, *&player.Curp)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "enfermedad = ?")
	if player.Enfermedad != nil {
		values = append(values, *&player.Enfermedad)
	} else {
		values = append(values, nil)
	}

	fields = append(fields, "insurance = ?")
	values = append(values, strconv.FormatBool(player.Insurance))

	fields = append(fields, "insurance_name = ?")
	if player.Insurance_name != nil {
		values = append(values, *&player.Insurance_name)
	} else {
		values = append(values, nil)
	}

	return fields, values
}
