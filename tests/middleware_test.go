package app_test

import (
	"checklist/middleware"
	"testing"
)

func TestUserCan(t *testing.T) {
	entityName1 := "users"
	entityName2 := "roles"
	entityName3 := "checklists"
	entityName4 := "checklist_items"
	action1 := "create"
	action2 := "read"
	action3 := "update"
	action4 := "delete"
	userDataValid := map[string]any{
		"abilities": `{
		    "roles": {
                "read": 1,
                "create": 1,
                "delete": 1,
                "update": 1
            },
            "users": {
                "read": 1,
                "create": 1,
                "delete": 1,
                "update": 1
            },
            "checklists": {
                "read": 1,
                "create": 1,
                "delete": 1,
                "update": 1
            },
            "checklist_items": {
                "read": 1,
                "create": 1,
                "delete": 1,
                "update": 1
            }
        }`,
	}
	if middleware.UserCan(userDataValid, entityName1, action1) {
		t.Log("Done with userDataValid for users")
	} else {
		t.Errorf("Result was incorrect with userDataValid for users")
	}
	if middleware.UserCan(userDataValid, entityName1, action2) {
		t.Log("Done with userDataValid for users")
	} else {
		t.Errorf("Result was incorrect with userDataValid for users")
	}
	if middleware.UserCan(userDataValid, entityName1, action3) {
		t.Log("Done with userDataValid for users")
	} else {
		t.Errorf("Result was incorrect with userDataValid for users")
	}
	if middleware.UserCan(userDataValid, entityName1, action4) {
		t.Log("Done with userDataValid for users")
	} else {
		t.Errorf("Result was incorrect with userDataValid for users")
	}

	if middleware.UserCan(userDataValid, entityName2, action1) {
		t.Log("Done with userDataValid for roles")
	} else {
		t.Errorf("Result was incorrect with userDataValid for roles")
	}
	if middleware.UserCan(userDataValid, entityName2, action2) {
		t.Log("Done with userDataValid for roles")
	} else {
		t.Errorf("Result was incorrect with userDataValid for roles")
	}
	if middleware.UserCan(userDataValid, entityName2, action3) {
		t.Log("Done with userDataValid for roles")
	} else {
		t.Errorf("Result was incorrect with userDataValid for roles")
	}
	if middleware.UserCan(userDataValid, entityName2, action4) {
		t.Log("Done with userDataValid for roles")
	} else {
		t.Errorf("Result was incorrect with userDataValid for roles")
	}

	if middleware.UserCan(userDataValid, entityName3, action1) {
		t.Log("Done with userDataValid for checklists")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklists")
	}
	if middleware.UserCan(userDataValid, entityName3, action2) {
		t.Log("Done with userDataValid for checklists")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklists")
	}
	if middleware.UserCan(userDataValid, entityName3, action3) {
		t.Log("Done with userDataValid for checklists")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklists")
	}
	if middleware.UserCan(userDataValid, entityName3, action4) {
		t.Log("Done with userDataValid for checklists")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklists")
	}

	if middleware.UserCan(userDataValid, entityName4, action1) {
		t.Log("Done with userDataValid for checklist_items")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklist_items")
	}
	if middleware.UserCan(userDataValid, entityName4, action2) {
		t.Log("Done with userDataValid for checklist_items")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklist_items")
	}
	if middleware.UserCan(userDataValid, entityName4, action3) {
		t.Log("Done with userDataValid for checklist_items")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklist_items")
	}
	if middleware.UserCan(userDataValid, entityName4, action4) {
		t.Log("Done with userDataValid for checklist_items")
	} else {
		t.Errorf("Result was incorrect with userDataValid for checklist_items")
	}
}
