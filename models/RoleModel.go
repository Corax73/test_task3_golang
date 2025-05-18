package models

type Role struct {
	*Model
}

func (role *Role) Init() *Role {
	model := Model{}
	model.SetTable("roles")
	model.Fields = map[string]string{"id": "", "title": "", "abilities": "", "created_at": ""}
	model.FieldTypes = map[string]string{"id": "int", "title": "string", "abilities": "string"}
	return &Role{&model}
}
