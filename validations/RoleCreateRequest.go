package validations

import (
	"checklist/customStructs"
	"fmt"
)

type RoleCreateValidatedFields struct {
	Title string
}
type RoleCreateValidatedData struct {
	Success bool
	Data    RoleCreateValidatedFields
}

func (roleCreateValidatedData *RoleCreateValidatedData) ToMap() map[string]any {
	resp := make(map[string]any, 1)
	resp["title"] = roleCreateValidatedData.Data.Title
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
	return response
}
