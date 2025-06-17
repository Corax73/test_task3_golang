package validations

import (
	"checklist/customStructs"
	"checklist/models"
	"checklist/utils"
	"fmt"
	"strconv"
)

type UserCreateValidatedFields struct {
	Login, Email, Password, RoleId, ChecklistsQuantity string
}
type UserCreateValidatedData struct {
	Success bool
	Data    UserCreateValidatedFields
}

func (userValidatedData *UserCreateValidatedData) ToMap() map[string]any {
	resp := make(map[string]any, 4)
	resp["login"] = userValidatedData.Data.Login
	resp["email"] = userValidatedData.Data.Email
	resp["password"] = userValidatedData.Data.Password
	resp["role_id"] = userValidatedData.Data.RoleId
	resp["checklists_quantity"] = userValidatedData.Data.ChecklistsQuantity
	return resp
}

func UserCreateRequestValidating(request customStructs.Request) UserCreateValidatedData {
	var response UserCreateValidatedData
	invalidData := "Invalid data"
	if login, ok := request.Params["login"]; ok && login != "" {
		response.Success = true
		response.Data.Login = fmt.Sprintf("%s", login)
	} else {
		response.Data.Login = invalidData
	}
	if email, ok := request.Params["email"]; ok && email != "" {
		emailStr := fmt.Sprintf("%s", email)
		if utils.IsEmail(emailStr) {
			response.Data.Email = emailStr
		} else {
			response.Success = false
			response.Data.Email = invalidData
		}
	} else {
		response.Success = false
		response.Data.Email = invalidData
	}
	if password, ok := request.Params["password"]; ok && password != "" {
		passwordStr := fmt.Sprintf("%s", password)
		if models.IsPasswordValid(passwordStr) {
			response.Data.Password = passwordStr
		} else {
			response.Success = false
			response.Data.Password = invalidData
		}
	} else {
		response.Success = false
		response.Data.Password = invalidData
	}
	if roleId, ok := request.Params["role_id"]; ok && roleId != "" {
		roleModel := (&models.Role{}).Init()
		roleIdInt := int(int64(roleId.(float64)))
		role := roleModel.GetOneById(roleIdInt)
		if role.Success {
			response.Data.RoleId = strconv.Itoa(roleIdInt)
		} else {
			response.Success = false
			response.Data.RoleId = invalidData
		}
	} else {
		response.Success = false
		response.Data.RoleId = invalidData
	}
	var checklistsQuantityInt int
	if checklistsQuantity, ok := request.Params["checklists_quantity"]; ok && checklistsQuantity != "" {
		checklistsQuantityInt = max(int(int64(checklistsQuantity.(float64))), 0)
	}
	response.Data.ChecklistsQuantity = strconv.Itoa(checklistsQuantityInt)
	return response
}
