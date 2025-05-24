package customStructs

import (
	"checklist/customLog"
	"encoding/json"
)

type SimpleResponse struct {
	Success bool
	Message map[string]any
}

type ListResponse struct {
	Success bool
	Message []map[string]any
}

func (listResponse *ListResponse) ToString() string {
	var resp string
	jsonData, err := json.Marshal(listResponse.Message)
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = string(jsonData)
	}
	return resp
}
