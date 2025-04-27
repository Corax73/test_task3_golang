package middleware

import (
	"checklist/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BasicCheck searches for a user using the passed login and compares the hash of the passed password with the password of the found user.
// Returns the result of the password check and user data
func BasicCheck(login, password string) (bool, map[string]any) {
	var isAuth bool
	userModel := (*&models.User{}).Init()
	res := userModel.GetOneByField("login", login, "roles")
	if res.Success {
		if err := bcrypt.CompareHashAndPassword([]byte(string(fmt.Sprintf("%s", res.Message["password"]))), []byte(password)); err == nil {
			isAuth = true
		}
	}
	return isAuth, res.Message
}
