package validations

import (
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

func RoleCreateRequestValidating(requestData map[string]any) RoleCreateValidatedData {
	var response RoleCreateValidatedData
	invalidData := "Invalid data"
	if title, ok := requestData["title"]; ok && title != "" {
		response.Success = true
		response.Data.Title = fmt.Sprintf("%s", title)
	} else {
		response.Data.Title = invalidData
	}
	return response
}
