package models

import (
	"checklist/customLog"
	"strconv"
)

type Checklist struct {
	*Model
}

func (checklist *Checklist) Init() *Checklist {
	model := Model{}
	model.SetTable("checklists")
	model.Fields = map[string]string{"id": "", "user_id": "", "title": "", "created_at": ""}
	model.FieldTypes = map[string]string{"id": "int", "user_id": "int", "checklist_id": "int", "title": "string", "created_at": "string"}
	return &Checklist{&model}
}

func (checklist *Checklist) CanCreating(userId string) bool {
	var resp bool
	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		customLog.Logging(err)
	} else {
		userModel := (&User{}).Init()
		user := userModel.GetOneById(intUserId)
		_, total := checklist.GetList(map[string]string{"filterBy": "user_id", "filterVal": userId})
		checklistsQuantity := user.Message["checklists_quantity"].(int64)
		if err != nil {
			customLog.Logging(err)
		}
		if checklistsQuantity > int64(total) {
			resp = true
		}
	}
	return resp
}
