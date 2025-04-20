package models

type Role struct {
	*Model
}

func (role *Role) Init() *Role {
	model := Model{}
	model.SetTable("roles")
	model.Fields = map[string]string{"id": "", "title": "", "created_at": ""}
	return &Role{&model}
}
