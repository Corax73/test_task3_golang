package models

import (
	"checklist/customDb"

	simplemodels "github.com/Corax73/simpleAbstractModels"
)

type ChecklistItem struct {
	*simplemodels.Model
}

func (role *ChecklistItem) Init() *ChecklistItem {
	model := simplemodels.Model{}
	model.SetDb(customDb.GetConnect())
	model.SetTable("checklist_items")
	model.Fields = map[string]string{"id": "", "is_completed": "0", "checklist_id": "", "description": "", "created_at": ""}
	model.FieldTypes = map[string]string{"id": "int", "is_completed": "int", "checklist_id": "int", "description": "string", "created_at": "string"}
	return &ChecklistItem{&model}
}
