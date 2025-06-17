package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/models"
	"strconv"
)

type EntityDeleteValidatedFields struct {
	Id string
}
type EntityDeleteValidatedData struct {
	Success bool
	Data    EntityDeleteValidatedFields
}

func EntityDeleteRequestValidating(request customStructs.Request, entityName string) EntityDeleteValidatedData {
	var response EntityDeleteValidatedData
	invalidData := "Invalid data"
	if id, ok := request.Params["id"]; ok && id != "" && entityName != "" {
		if entityName == "users" {
			userModel := (&models.User{}).Init()
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
		} else if entityName == "roles" {
			roleModel := (&models.Role{}).Init()
			roleIdInt, err := strconv.Atoi(id.(string))
			if err != nil {
				customLog.Logging(err)
				response.Success = false
				response.Data.Id = invalidData
			} else {
				role := roleModel.GetOneById(roleIdInt)
				if role.Success {
					response.Data.Id = strconv.Itoa(roleIdInt)
					response.Success = true
				} else {
					response.Success = false
					response.Data.Id = invalidData
				}
			}
		} else if entityName == "checklists" {
			checklistModel := (&models.Checklist{}).Init()
			checklistIdInt, err := strconv.Atoi(id.(string))
			if err != nil {
				customLog.Logging(err)
				response.Success = false
				response.Data.Id = invalidData
			} else {
				checklist := checklistModel.GetOneById(checklistIdInt)
				if checklist.Success {
					response.Data.Id = strconv.Itoa(checklistIdInt)
					response.Success = true
				} else {
					response.Success = false
					response.Data.Id = invalidData
				}
			}
		} else if entityName == "checklist_items" {
			checklistItemModel := (&models.ChecklistItem{}).Init()
			checklistItemIdInt, err := strconv.Atoi(id.(string))
			if err != nil {
				customLog.Logging(err)
				response.Success = false
				response.Data.Id = invalidData
			} else {
				checklist := checklistItemModel.GetOneById(checklistItemIdInt)
				if checklist.Success {
					response.Data.Id = strconv.Itoa(checklistItemIdInt)
					response.Success = true
				} else {
					response.Success = false
					response.Data.Id = invalidData
				}
			}
		}
	} else {
		response.Success = false
		response.Data.Id = invalidData
	}
	return response
}
