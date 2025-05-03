package validations

import (
	"checklist/customStructs"
	"checklist/models"
	"fmt"
	"slices"
	"strconv"
)

type ChecklistItemCreateValidatedFields struct {
	ChecklistCreateValidatedData
	ChecklistId, IsCompleted, Description string
}
type ChecklistItemCreateValidatedData struct {
	Success bool
	Data    ChecklistItemCreateValidatedFields
}

func (checklistItemCreateValidatedData *ChecklistItemCreateValidatedData) ToMap() map[string]any {
	resp := make(map[string]any, 2)
	resp["checklist_id"] = checklistItemCreateValidatedData.Data.ChecklistId
	resp["is_completed"] = checklistItemCreateValidatedData.Data.IsCompleted
	resp["description"] = checklistItemCreateValidatedData.Data.Description
	return resp
}

func ChecklistItemCreateRequestValidating(request customStructs.Request) ChecklistItemCreateValidatedData {
	var response ChecklistItemCreateValidatedData
	invalidData := "Invalid data"
	response.Data.IsCompleted = invalidData
	if isCompleted, ok := request.Params["is_completed"]; ok && isCompleted != "" {
		strVal := fmt.Sprintf("%v", isCompleted)
		if slices.Contains([]string{"0", "1"}, strVal) {
			response.Success = true
			response.Data.IsCompleted = strVal
		}
	}
	if description, ok := request.Params["description"]; ok && description != "" {
		response.Data.Description = fmt.Sprintf("%s", description)
	} else {
		response.Success = false
		response.Data.Description = invalidData
	}
	if checklistId, ok := request.Params["checklist_id"]; ok && checklistId != "" {
		checklistModel := (*&models.Checklist{}).Init()
		checklistIdInt := int(int64(checklistId.(float64)))
		checklist := checklistModel.GetOneById(checklistIdInt)
		if checklist.Success {
			response.Data.ChecklistId = strconv.Itoa(checklistIdInt)
		} else {
			response.Success = false
			response.Data.ChecklistId = invalidData
		}
	} else {
		response.Success = false
		response.Data.ChecklistId = invalidData
	}
	return response
}
