package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/models"
	"checklist/utils"
	"fmt"
	"strconv"
)

type UserUpdateValidatedFields struct {
	Id, Login, Email, Password, RoleId string
}
type UserUpdateValidatedData struct {
	Success bool
	Data    UserUpdateValidatedFields
}

func (userValidatedData *UserUpdateValidatedData) ToMap() map[string]string {
	resp := make(map[string]string, 5)
	resp["id"] = userValidatedData.Data.Id
	resp["login"] = userValidatedData.Data.Login
	resp["email"] = userValidatedData.Data.Email
	resp["password"] = userValidatedData.Data.Password
	resp["role_id"] = userValidatedData.Data.RoleId
	return resp
}

func UserUpdateRequestValidating(request customStructs.Request) UserUpdateValidatedData {
	var response UserUpdateValidatedData
	invalidData := "Invalid data"
	if login, ok := request.Params["login"]; ok && login != "" {
		response.Success = true
		response.Data.Login = fmt.Sprintf("%s", login)
	}
	if email, ok := request.Params["email"]; ok && email != "" {
		emailStr := fmt.Sprintf("%s", email)
		if utils.IsEmail(emailStr) {
			response.Data.Email = emailStr
		} else {
			response.Success = false
			response.Data.Email = invalidData
		}
	}
	if password, ok := request.Params["password"]; ok && password != "" {
		passwordStr := fmt.Sprintf("%s", password)
		if models.IsPasswordValid(passwordStr) {
			response.Data.Password = passwordStr
		} else {
			response.Success = false
			response.Data.Password = invalidData
		}
	}
	if roleId, ok := request.Params["role_id"]; ok && roleId != "" {
		roleModel := (*&models.Role{}).Init()
		roleIdInt := int(int64(roleId.(float64)))
		role := roleModel.GetOneById(roleIdInt)
		if role.Success {
			response.Data.RoleId = strconv.Itoa(roleIdInt)
		} else {
			response.Success = false
			response.Data.RoleId = invalidData
		}
	}
	if id, ok := request.Params["id"]; ok && id != "" {
		userModel := (*&models.User{}).Init()
		userIdInt, err := strconv.Atoi(id.(string))
		if err != nil {
			customLog.Logging(err)
			response.Success = false
			response.Data.Id = invalidData
		} else {
			user := userModel.GetOneById(userIdInt)
			if user.Success {
				response.Data.Id = strconv.Itoa(userIdInt)
			} else {
				response.Success = false
				response.Data.Id = invalidData
			}
		}
	} else {
		response.Success = false
		response.Data.Id = invalidData
	}
	return response
}
