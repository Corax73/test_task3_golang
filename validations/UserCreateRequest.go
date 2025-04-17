package validations

import (
	"checklist/models"
	"checklist/utils"
	"fmt"
	"strconv"
)

type UserCreateValidatedFields struct {
	Login, Email, Password, RoleId string
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
	return resp
}

func UserCreateRequestValidating(requestData map[string]any) UserCreateValidatedData {
	var response UserCreateValidatedData
	invalidData := "Invalid data"
	if login, ok := requestData["login"]; ok && login != "" {
		response.Success = true
		response.Data.Login = fmt.Sprintf("%s", login)
	} else {
		response.Data.Login = invalidData
	}
	if email, ok := requestData["email"]; ok && email != "" {
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
	if password, ok := requestData["password"]; ok && password != "" {
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
	if roleId, ok := requestData["role_id"]; ok && roleId != "" {
		roleModel := (*&models.Role{}).Init()
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
	return response
}
