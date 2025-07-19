package validations

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/utils"
	"encoding/json"
	"fmt"
	"strings"
)

type EntityListValidatedFields struct {
	Id, FilterBy, FilterVal, OrderBy, Order, Limit, Offset string
	Filters                                                []map[string]string
}
type EntityListValidatedData struct {
	Success bool
	Data    EntityListValidatedFields
}

func (entityValidatedData *EntityListValidatedData) ToMap() map[string]string {
	resp := make(map[string]string, 7)
	resp["id"] = entityValidatedData.Data.Id
	resp["filterBy"] = entityValidatedData.Data.FilterBy
	resp["filterVal"] = entityValidatedData.Data.FilterVal
	resp["orderBy"] = entityValidatedData.Data.OrderBy
	resp["order"] = entityValidatedData.Data.Order
	resp["limit"] = entityValidatedData.Data.Limit
	resp["offset"] = entityValidatedData.Data.Offset
	return resp
}

func EntityListRequestValidating(request customStructs.Request) EntityListValidatedData {
	var response EntityListValidatedData
	if id, ok := request.Params["id"]; ok && id != "" {
		response.Data.Id = fmt.Sprintf("%s", id)
		response.Success = true
	}
	if filterBy, ok := request.Params["filterBy"]; ok && filterBy != "" {
		response.Data.FilterBy = fmt.Sprintf("%s", filterBy)
		if filterVal, ok := request.Params["filterVal"]; ok && filterVal != "" {
			response.Success = true
			response.Data.FilterVal = fmt.Sprintf("%s", filterVal)
		}
	}
	if orderBy, ok := request.Params["orderBy"]; ok && orderBy != "" {
		response.Success = true
		response.Data.OrderBy = fmt.Sprintf("%s", orderBy)
		if order, ok := request.Params["order"]; ok && order != "" {
			response.Success = true
			response.Data.Order = fmt.Sprintf("%s", order)
		}
	}
	if limit, ok := request.Params["limit"]; ok && limit != "" {
		response.Success = true
		response.Data.Limit = fmt.Sprintf("%s", limit)
	}
	if offset, ok := request.Params["offset"]; ok && offset != "" {
		response.Success = true
		response.Data.Offset = fmt.Sprintf("%s", offset)
	}
	if len(request.Filters) > 0 {
		response.Data.Filters = request.Filters
	}
	return response
}

func (entityValidatedData *EntityListValidatedData) GetAsKey(entityName string) string {
	jsonString, err := json.Marshal(entityValidatedData.Data.Filters)
	if err != nil {
		customLog.Logging(err)
	}
	filterKey := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(
							strings.ReplaceAll(
								string(jsonString), "[", "-",
							),
							"\"", "-",
						), ",", "-",
					), "{", "-"),
				":", "-"),
			"}", "-"),
		"]", "-")
	return utils.ConcatSlice([]string{
		entityName,
		"-",
		entityValidatedData.Data.Id,
		"-",
		entityValidatedData.Data.FilterBy,
		"-",
		entityValidatedData.Data.FilterVal,
		"-",
		entityValidatedData.Data.OrderBy,
		"-",
		entityValidatedData.Data.Order,
		"-",
		entityValidatedData.Data.Limit,
		"-",
		entityValidatedData.Data.Offset,
		"-",
		filterKey,
	})
}
