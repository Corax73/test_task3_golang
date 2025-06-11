package utilstest

import (
	"checklist/utils"
	"testing"
)

func TestGetConfFromEnvFile(t *testing.T) {
	envData := utils.GetConfFromEnvFile("1")
	if len(envData) == 0 {
		t.Log("Done with incorrect filename")
	} else {
		t.Errorf("Result was incorrect with incorrect filename")
	}
	envData = utils.GetConfFromEnvFile("./env.test")
	if len(envData) > 0 {
		t.Log("Done with correct filename")
	} else {
		t.Errorf("Result was incorrect with correct filename")
	}
	envData = utils.GetConfFromEnvFile("")
	if len(envData) > 0 {
		t.Errorf("Result was incorrect with correct filename")
	} else {
		t.Log("Done with incorrect filename")
	}
}