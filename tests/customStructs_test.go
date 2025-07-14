package app_test

import (
	"checklist/customStructs"
	"testing"
)

func TestToString(t *testing.T) {
	simpleResponse := customStructs.ListResponse{
		Success: true,
		Message: []map[string]any{
			{
				"key1": 1,
				"key2": "2",
				"key3": true,
			},
			{
				"key1": 1,
				"key2": 2,
				"key3": []int{1, 2, 3},
			},
			{
				"key1": map[string]int{"1": 1, "2": 2, "3": 3},
			},
		},
		Total: 3,
	}
	validString := `{"data":[{"key1":1,"key2":"2","key3":true},{"key1":1,"key2":2,"key3":[1,2,3]},{"key1":{"1":1,"2":2,"3":3}}],"total":3}`
	inValidString := `[{"key1":1,"key2":"2","key3":true},{"key1":1,"key2":2,"key3":[1,2,3]}]`
	if simpleResponse.ToString() == validString {
		t.Log("Done with validString")
	} else {
		t.Errorf("Result was incorrect with validString")
	}
	if simpleResponse.ToString() == inValidString {
		t.Errorf("Result was incorrect with inValidString")
	} else {
		t.Log("Done with inValidString")
	}
}
