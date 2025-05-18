package models

type Checklist struct {
	*Model
}

func (role *Checklist) Init() *Checklist {
	model := Model{}
	model.SetTable("checklists")
	model.Fields = map[string]string{"id": "", "user_id": "", "title": "", "created_at": ""}
	model.FieldTypes = map[string]string{"id": "int", "user_id": "int", "checklist_id": "int", "title": "string"}
	return &Checklist{&model}
}
