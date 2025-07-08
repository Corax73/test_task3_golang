package models

type ChecklistItem struct {
	*Model
}

func (role *ChecklistItem) Init() *ChecklistItem {
	model := Model{}
	model.SetTable("checklist_items")
	model.Fields = map[string]string{"id": "", "is_completed": "0", "checklist_id": "", "description": "", "created_at": ""}
	model.FieldTypes = map[string]string{"id": "int", "is_completed": "int", "checklist_id": "int", "description": "string", "created_at": "string"}
	return &ChecklistItem{&model}
}
