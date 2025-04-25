package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"encoding/json"
	"fmt"
)

type RoleCreateValidatedFields struct {
	Title, Abilities string
}
type RoleCreateValidatedData struct {
	Success bool
	Data    RoleCreateValidatedFields
}

func (roleCreateValidatedData *RoleCreateValidatedData) ToMap() map[string]any {
	resp := make(map[string]any, 2)
	resp["title"] = roleCreateValidatedData.Data.Title
	resp["abilities"] = roleCreateValidatedData.Data.Abilities
	return resp
}

func RoleCreateRequestValidating(request customStructs.Request) RoleCreateValidatedData {
	var response RoleCreateValidatedData
	invalidData := "Invalid data"
	if title, ok := request.Params["title"]; ok && title != "" {
		response.Success = true
		response.Data.Title = fmt.Sprintf("%s", title)
	} else {
		response.Data.Title = invalidData
	}
	if abilities, ok := request.Params["abilities"]; ok {
		bs, err := json.Marshal(abilities)
		if err != nil {
			customLog.Logging(err)
		} else {
			response.Data.Abilities = string(bs)
		}
	}
	return response
}
