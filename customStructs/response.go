package customStructs

import (
	"checklist/customLog"
	"checklist/utils"
	"encoding/json"
	"strconv"
)

type SimpleResponse struct {
	Success bool
	Message map[string]any
}

type ListResponse struct {
	Success bool
	Message []map[string]any `json:"data"`
	Total   int `json:"total"`
}

func (listResponse *ListResponse) ToString() string {
	var resp string
	jsonData, err := json.Marshal(listResponse.Message)
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = string(jsonData)
		resp = utils.ConcatSlice([]string{ "{\"data\":", resp, ",", "\"total\":", strconv.Itoa(listResponse.Total), "}"})
	}
	return resp
}
