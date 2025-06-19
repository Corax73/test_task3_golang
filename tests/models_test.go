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

func TestSetTable(t *testing.T) {
	user := (&models.User{}).Init()
	user.SetTable("users1")
	if user.Table() == "users1" {
		t.Log("Done for model User")
	} else {
		t.Errorf("Result was incorrect for model User")
	}
	role := (&models.Role{}).Init()
	role.SetTable("roles1")
	if role.Table() == "roles1" {
		t.Log("Done for model Role")
	} else {
		t.Errorf("Result was incorrect for model Role")
	}
	checklist := (&models.Checklist{}).Init()
	checklist.SetTable("checklists1")
	if checklist.Table() == "checklists1" {
		t.Log("Done for model Checklist")
	} else {
		t.Errorf("Result was incorrect for model Checklist")
	}
	checklistItem := (&models.ChecklistItem{}).Init()
	checklistItem.SetTable("checklist_items1")
	if checklistItem.Table() == "checklist_items1" {
		t.Log("Done for model ChecklistItem")
	} else {
		t.Errorf("Result was incorrect for model ChecklistItem")
	}
}

func TestInit(t *testing.T) {
	userFieldsMap := map[string]string{"id": "", "role_id": "", "login": "", "email": "", "password": "", "checklists_quantity": "0", "created_at": ""}
	user := (&models.User{}).Init()
	for key, value := range user.Fields {
		if v, ok := userFieldsMap[key]; !ok {
			t.Errorf("Keys was incorrect for model User")
		} else {
			if value == v {
				t.Log("Done for model User")
			} else {
				t.Errorf("Values was incorrect for model User")
			}
		}
	}
	roleFieldsMap := map[string]string{"id": "", "title": "", "abilities": "", "created_at": ""}
	role := (&models.Role{}).Init()
	for key, value := range role.Fields {
		if v, ok := roleFieldsMap[key]; !ok {
			t.Errorf("Keys was incorrect for model Role")
		} else {
			if value == v {
				t.Log("Done for model Role")
			} else {
				t.Errorf("Values was incorrect for model Role")
			}
		}
	}
	checklistFieldsMap := map[string]string{"id": "", "user_id": "", "title": "", "created_at": ""}
	checklist := (&models.Checklist{}).Init()
	for key, value := range checklist.Fields {
		if v, ok := checklistFieldsMap[key]; !ok {
			t.Errorf("Keys was incorrect for model Checklist")
		} else {
			if value == v {
				t.Log("Done for model Checklist")
			} else {
				t.Errorf("Values was incorrect for model Checklist")
			}
		}
	}

	checklistItemFieldsMap := map[string]string{"id": "", "is_completed": "0", "checklist_id": "", "description": "", "created_at": ""}
	checklistItem := (&models.ChecklistItem{}).Init()
	for key, value := range checklistItem.Fields {
		if v, ok := checklistItemFieldsMap[key]; !ok {
			t.Errorf("Keys was incorrect for model ChecklistItem")
		} else {
			if value == v {
				t.Log("Done for model ChecklistItem")
			} else {
				t.Errorf("Values was incorrect for model ChecklistItem")
			}
		}
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
