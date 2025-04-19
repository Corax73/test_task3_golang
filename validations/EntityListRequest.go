package validations

import (
	"fmt"
)

type EntityListValidatedFields struct {
	FilterBy, FilterVal, OrderBy, Order, Limit, Offset string
}
type EntityListValidatedData struct {
	Success bool
	Data    EntityListValidatedFields
}

func (entityValidatedData *EntityListValidatedData) ToMap() map[string]string {
	resp := make(map[string]string, 6)
	resp["filterBy"] = entityValidatedData.Data.FilterBy
	resp["filterVal"] = entityValidatedData.Data.FilterVal
	resp["orderBy"] = entityValidatedData.Data.OrderBy
	resp["order"] = entityValidatedData.Data.Order
	resp["limit"] = entityValidatedData.Data.Limit
	resp["offset"] = entityValidatedData.Data.Offset
	return resp
}

func EntityListRequestValidating(requestData map[string]any) EntityListValidatedData {
	var response EntityListValidatedData
	if filterBy, ok := requestData["filterBy"]; ok && filterBy != "" {
		response.Data.FilterBy = fmt.Sprintf("%s", filterBy)
		if filterVal, ok := requestData["filterVal"]; ok && filterVal != "" {
			response.Success = true
			response.Data.FilterVal = fmt.Sprintf("%s", filterVal)
		}
	}
	if orderBy, ok := requestData["orderBy"]; ok && orderBy != "" {
		response.Success = true
		response.Data.OrderBy = fmt.Sprintf("%s", orderBy)
		if order, ok := requestData["order"]; ok && order != "" {
			response.Success = true
			response.Data.Order = fmt.Sprintf("%s", order)
		}
	}
	if limit, ok := requestData["limit"]; ok && limit != "" {
		response.Success = true
		response.Data.Limit = fmt.Sprintf("%s", limit)
	}
	if offset, ok := requestData["offset"]; ok && offset != "" {
		response.Success = true
		response.Data.Offset = fmt.Sprintf("%s", offset)
	}
	return response
}
