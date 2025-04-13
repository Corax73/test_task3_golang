package models

type User struct {
	*Model
}

func (song *User) Init() *User {
	model := Model{}
	model.SetTable("users")
	model.Fields = map[string]string{"id": "", "role_id": "", "login": "", "email": "", "password": "", "created_at": ""}
	return &User{&model}
}
