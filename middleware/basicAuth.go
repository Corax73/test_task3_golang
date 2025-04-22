package middleware

import (
	"checklist/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BasicCheck searches for a user using the passed login and compares the hash of the passed password with the password of the found user
func BasicCheck(login, password string) bool {
	var resp bool
	userModel := (*&models.User{}).Init()
	res := userModel.GetOneByField("login", login)
	if res.Success {
		if err := bcrypt.CompareHashAndPassword([]byte(string(fmt.Sprintf("%s", res.Message["password"]))), []byte(password)); err == nil {
			resp = true
		}
	}
	return resp
}
