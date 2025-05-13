package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/models"
	"fmt"
	"strconv"
)

type ChecklistUpdateValidatedFields struct {
	Id, UserId, Title string
}
type ChecklistUpdateValidatedData struct {
	Success bool
	Data    ChecklistUpdateValidatedFields
}

func (checklistValidatedData *ChecklistUpdateValidatedData) ToMap() map[string]string {
	resp := make(map[string]string, 3)
	resp["id"] = checklistValidatedData.Data.Id
	resp["user_id"] = checklistValidatedData.Data.UserId
	resp["title"] = checklistValidatedData.Data.Title
	return resp
}

func ChecklistUpdateRequestValidating(request customStructs.Request) ChecklistUpdateValidatedData {
	var response ChecklistUpdateValidatedData
	invalidData := "Invalid data"
	if title, ok := request.Params["title"]; ok && title != "" {
		response.Data.Title = fmt.Sprintf("%s", title)
	} else {
		response.Data.Title = invalidData
	}
	if userId, ok := request.Params["user_id"]; ok && userId != "" {
		userModel := (*&models.User{}).Init()
		userIdInt := int(int64(userId.(float64)))
		user := userModel.GetOneById(userIdInt)
		if user.Success {
			response.Data.UserId = strconv.Itoa(userIdInt)
		} else {
			response.Success = false
			response.Data.UserId = invalidData
		}
	}
	if id, ok := request.Params["id"]; ok && id != "" {
		userModel := (*&models.Checklist{}).Init()
		userIdInt, err := strconv.Atoi(id.(string))
		if err != nil {
			customLog.Logging(err)
			response.Success = false
			response.Data.Id = invalidData
		} else {
			user := userModel.GetOneById(userIdInt)
			if user.Success {
				response.Data.Id = strconv.Itoa(userIdInt)
				response.Data.Id = strconv.Itoa(userIdInt)
				if response.Data.Title != invalidData ||
					response.Data.UserId != invalidData {
					response.Success = true
				}
				if response.Data.Title == invalidData ||
					response.Data.UserId == invalidData {
					response.Success = false
				}
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
