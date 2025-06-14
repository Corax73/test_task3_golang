package utilstest

import (
	"checklist/utils"
	"slices"
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

func TestGetMapKeysWithValue(t *testing.T) {
	map1 := map[string]string{"1": "1", "2": "2"}
	map2 := map[string]string{"1": "", "2": "2"}
	if len(utils.GetMapKeysWithValue(map1)) == 2 {
		t.Log("Done with a map with two non-empty values")
	} else {
		t.Errorf("Result was incorrect with a map with two non-empty values")
	}
	if len(utils.GetMapKeysWithValue(map2)) == 1 {
		t.Log("Done with a map with one non-empty values")
	} else {
		t.Errorf("Result was incorrect with a map with one non-empty values")
	}
}

func TestGetMapValues(t *testing.T) {
	map1 := map[string]string{"1": "1", "2": "2"}
	map2 := map[string]string{"1": "", "2": "2"}
	testSlice1 := utils.GetMapValues(map1)
	var hasErr bool
	for _, val := range map1 {
		if !slices.Contains(testSlice1, val) {
			t.Errorf("Result was incorrect for value %s of map1", val)
			hasErr = true
		}
	}
	if !hasErr {
		t.Log("Done with a map with map1")
	}
	testSlice2 := utils.GetMapValues(map2)
	hasErr = false
	for _, val := range map2 {
		if !slices.Contains(testSlice2, val) && val != "" {
			t.Errorf("Result was incorrect for value %s of map2", val)
			hasErr = true
		}
	}
	if !hasErr {
		t.Log("Done with a map with map2")
	}
}

func TestGetIndexByStrValue(t *testing.T) {
	testSlice := []string{"1", "2"}
	if utils.GetIndexByStrValue(testSlice, "2") == 1 {
		t.Log("Done with a map with correct index")
	} else {
		t.Errorf("Result was incorrect with correct index")
	}
	if utils.GetIndexByStrValue(testSlice, "2") == 0 {
		t.Errorf("Result was incorrect with incorrect index")
	} else {
		t.Log("Done with a map with incorrect index")
	}
}

func TestGetMapKeys(t *testing.T) {
	map1 := map[string]string{"1": "1", "2": "2"}
	keys1 := []string{"1", "2"}
	map2 := map[string]string{"3": "3", "4": "4"}
	keys2 := []string{"3", "4"}
	var hasErr bool
	for key, val := range map1 {
		if !slices.Contains(keys1, key) {
			t.Errorf("Result was incorrect for key %s of map1", val)
			hasErr = true
		}
	}
	if !hasErr {
		t.Log("Done with a map with map1")
	}
	hasErr = false
	for key, val := range map2 {
		if !slices.Contains(keys2, key) {
			t.Errorf("Result was incorrect for key %s of map2", val)
			hasErr = true
		}
	}
	if !hasErr {
		t.Log("Done with a map with map2")
	}
}
