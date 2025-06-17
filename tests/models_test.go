package app_test

import (
	"checklist/models"
	"testing"
)

func TestTable(t *testing.T) {
	user := (&models.User{}).Init()
	if user.Table() == "users" {
		t.Log("Done for model User")
	} else {
		t.Errorf("Result was incorrect for model User")
	}
	role := (&models.Role{}).Init()
	if role.Table() == "roles" {
		t.Log("Done for model Role")
	} else {
		t.Errorf("Result was incorrect for model Role")
	}
	checklist := (&models.Checklist{}).Init()
	if checklist.Table() == "checklists" {
		t.Log("Done for model Checklist")
	} else {
		t.Errorf("Result was incorrect for model Checklist")
	}
	checklistItem := (&models.ChecklistItem{}).Init()
	if checklistItem.Table() == "checklist_items" {
		t.Log("Done for model ChecklistItem")
	} else {
		t.Errorf("Result was incorrect for model ChecklistItem")
	}
}

func TestIsPasswordValid(t *testing.T) {
	validPassword := "12345_Dd"
	invalidPassword1 := "1234_Dd"
	invalidPassword2 := "12345_DФ"
	if models.IsPasswordValid(validPassword) {
		t.Log("Done with valid password")
	} else {
		t.Errorf("Result was incorrect with valid password")
	}
	if models.IsPasswordValid(invalidPassword1) {
		t.Errorf("Result was incorrect with invalidPassword1")
	} else {
		t.Log("Done with invalidPassword1")
	}
	if models.IsPasswordValid(invalidPassword2) {
		t.Errorf("Result was incorrect with invalidPassword2")
	} else {
		t.Log("Done with invalidPassword2")
	}
}
