package models

import (
	"checklist/customDb"

	simplemodels "github.com/Corax73/simpleAbstractModels"
)

type Role struct {
	*simplemodels.Model
}

func (role *Role) Init() *Role {
	model := simplemodels.Model{}
	model.SetDb(customDb.GetConnect())
	model.SetTable("roles")
	model.Fields = map[string]string{"id": "", "title": "", "abilities": "", "created_at": ""}
	model.FieldTypes = map[string]string{"id": "int", "title": "string", "abilities": "string", "created_at": "string"}
	return &Role{&model}
}
