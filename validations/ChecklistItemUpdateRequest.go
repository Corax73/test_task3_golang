package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/models"
	"fmt"
	"slices"
	"strconv"
)

type ChecklistItemUpdateValidatedFields struct {
	Id, IsCompleted, Description string
}
type ChecklistItemUpdateValidatedData struct {
	Success bool
	Data    ChecklistItemUpdateValidatedFields
}

func (checklistItemValidatedData *ChecklistItemUpdateValidatedData) ToMap() map[string]string {
	resp := make(map[string]string, 2)
	resp["id"] = checklistItemValidatedData.Data.Id
	resp["is_completed"] = checklistItemValidatedData.Data.IsCompleted
	resp["description"] = checklistItemValidatedData.Data.Description
	return resp
}

func ChecklistItemUpdateRequestValidating(request customStructs.Request) ChecklistItemUpdateValidatedData {
	var response ChecklistItemUpdateValidatedData
	invalidData := "Invalid data"
	response.Data.IsCompleted = invalidData
	if isCompleted, ok := request.Params["is_completed"]; ok && isCompleted != "" {
		strVal := fmt.Sprintf("%v", isCompleted)
		if slices.Contains([]string{"0", "1"}, strVal) {
			response.Data.IsCompleted = strVal
		} else {
			response.Success = false
			response.Data.IsCompleted = invalidData
		}
	}
	if description, ok := request.Params["description"]; ok && description != "" {
		response.Data.Description = fmt.Sprintf("%s", description)
	}
	if id, ok := request.Params["id"]; ok && id != "" {
		checklistItemModel := (&models.ChecklistItem{}).Init()
		checklistItemIdInt, err := strconv.Atoi(id.(string))
		if err != nil {
			customLog.Logging(err)
			response.Success = false
			response.Data.Id = invalidData
		} else {
			checklistItem := checklistItemModel.GetOneById(checklistItemIdInt)
			if checklistItem.Success {
				response.Data.Id = strconv.Itoa(checklistItemIdInt)
				if response.Data.IsCompleted != invalidData ||
					response.Data.Description != invalidData {
					response.Success = true
				}
				if response.Data.IsCompleted == invalidData ||
					response.Data.Description == invalidData {
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
