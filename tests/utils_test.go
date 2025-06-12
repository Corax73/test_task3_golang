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

func TestConcatSlice(t *testing.T) {
	errStr := "321"
	correctStr := "123"
	if errStr != utils.ConcatSlice([]string{"1", "2", "3"}) {
		t.Log("Done with incorrect comparison")
	} else {
		t.Errorf("Result was incorrect with incorrect comparison")
	}
	if correctStr == utils.ConcatSlice([]string{"1", "2", "3"}) {
		t.Log("Done with correct comparison")
	} else {
		t.Errorf("Result was correct with incorrect comparison")
	}
}

func TestCompareMapsByStringKeys(t *testing.T) {
	map1 := map[string]string{"1": "1", "2": "2"}
	map2 := map[string]string{"1": "1", "2": "2"}
	map3 := map[string]string{"1": "1", "2": "2", "3": "3"}
	if utils.CompareMapsByStringKeys(map1, map2) {
		t.Log("Done with identical maps")
	} else {
		t.Errorf("Result was incorrect with identical maps")
	}
	if !utils.CompareMapsByStringKeys(map1, map3) {
		t.Log("Done with unequal maps")
	} else {
		t.Errorf("Result was incorrect with unequal maps")
	}
}
