package validations

import (
	"checklist/models"
	"fmt"
	"strconv"
)

type ChecklistCreateValidatedFields struct {
	UserId, Title string
}
type ChecklistCreateValidatedData struct {
	Success bool
	Data    ChecklistCreateValidatedFields
}

func (checklistValidatedData *ChecklistCreateValidatedData) ToMap() map[string]any {
	resp := make(map[string]any, 2)
	resp["user_id"] = checklistValidatedData.Data.UserId
	resp["title"] = checklistValidatedData.Data.Title
	return resp
}

func ChecklistCreateRequestValidating(requestData map[string]any) ChecklistCreateValidatedData {
	var response ChecklistCreateValidatedData
	invalidData := "Invalid data"
	if title, ok := requestData["title"]; ok && title != "" {
		response.Success = true
		response.Data.Title = fmt.Sprintf("%s", title)
	} else {
		response.Data.Title = invalidData
	}
	if userId, ok := requestData["user_id"]; ok && userId != "" {
		userModel := (*&models.User{}).Init()
		userIdInt := int(int64(userId.(float64)))
		role := userModel.GetOneById(userIdInt)
		if role.Success {
			response.Data.UserId = strconv.Itoa(userIdInt)
		} else {
			response.Success = false
			response.Data.UserId = invalidData
		}
	} else {
		response.Success = false
		response.Data.UserId = invalidData
	}
	return response
}
