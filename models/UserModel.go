package models

import (
	"regexp"
	"unicode/utf8"
)

type User struct {
	*Model
}

func (user *User) Init() *User {
	model := Model{}
	model.SetTable("users")
	model.Fields = map[string]string{"id": "", "role_id": "", "login": "", "email": "", "password": "", "created_at": ""}
	return &User{&model}
}

func (user *User) IsPasswordValid(password string) bool {
	lowerCond := regexp.MustCompile(`[a-z]`)
	upperCond := regexp.MustCompile(`[A-Z]`)
	digitCond := regexp.MustCompile(`[0-9]`)
	wholeCond := regexp.MustCompile(`^[0-9a-zA-Z!_@#$%^&*()]*$`)
	passLen := utf8.RuneCountInString(password) >= 8
	return lowerCond.MatchString(password) &&
		upperCond.MatchString(password) &&
		digitCond.MatchString(password) &&
		wholeCond.MatchString(password) &&
		passLen
}
