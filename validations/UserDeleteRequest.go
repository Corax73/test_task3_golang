package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/models"
	"strconv"
)

type UserDeleteValidatedFields struct {
	Id string
}
type UserDeleteValidatedData struct {
	Success bool
	Data    UserDeleteValidatedFields
}

func UserDeleteRequestValidating(request customStructs.Request) UserDeleteValidatedData {
	var response UserDeleteValidatedData
	invalidData := "Invalid data"
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
				response.Success = true
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
