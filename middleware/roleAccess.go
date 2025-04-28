package middleware

import "encoding/json"

// UserCan checks for the presence of the passed string key EntityName in the passed map and the presence of the string key Action.
// If they are present, returns a comparison of the value under the key Action and one, otherwise false
func UserCan(userData map[string]any, entityName, action string) bool {
	var resp bool
	abilities := make(map[string]map[string]int, 4)
	json.Unmarshal([]byte(userData["abilities"].(string)), &abilities)
	if accessData, ok := abilities[entityName]; ok {
		if actionAccess, ok := accessData[action]; ok {
			if actionAccess == 1 {
				resp = true
			}
		}
	}
	return resp
}
