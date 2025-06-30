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
	Message []map[string]any
	Total   int
}

func (listResponse *ListResponse) ToString() string {
	var resp string
	jsonData, err := json.Marshal(listResponse.Message)
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = string(jsonData)
		resp = utils.ConcatSlice([]string{resp, "\"total\":", strconv.Itoa(listResponse.Total)})
	}
	return resp
}
