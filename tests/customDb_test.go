package app_test

import (
	"checklist/customDb"
	"testing"
)

func TestGetDsnString(t *testing.T) {
	validEnvMap := map[string]string{
		"DB_USER":     "user",
		"DB_PASSWORD": "12345678",
		"DB_NAME":     "name",
		"DB_SSLMODE":  "disable",
	}
	invalidEnvMap := map[string]string{
		"1":     "user",
		"DB_PASSWORD": "12345678",
		"DB_NAME":     "name",
		"DB_SSLMODE":  "disable",
	}
	envStr := customDb.GetDsnString(validEnvMap)
	if envStr == "user=user password=12345678 dbname=name sslmode=disable " {
		t.Log("Done with validEnvMap")
	} else {
		t.Errorf("Result was incorrect with validEnvMap")
	}
	envStr = customDb.GetDsnString(invalidEnvMap)
	if envStr == "user=user password=12345678 dbname=name sslmode=disable " {
		t.Errorf("Result was incorrect with invalidEnvMap")
	} else {
		t.Log("Done with invalidEnvMap")
	}
}
