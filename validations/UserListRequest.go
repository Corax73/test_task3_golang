package validations

import (
	"fmt"
)

type UserListValidatedFields struct {
	FilterBy, FilterVal, OrderBy, Order, Limit, Offset string
}
type UserListValidatedData struct {
	Success bool
	Data    UserListValidatedFields
}

func (userValidatedData *UserListValidatedData) ToMap() map[string]string {
	resp := make(map[string]string, 6)
	resp["filterBy"] = userValidatedData.Data.FilterBy
	resp["filterVal"] = userValidatedData.Data.FilterVal
	resp["orderBy"] = userValidatedData.Data.OrderBy
	resp["order"] = userValidatedData.Data.Order
	resp["limit"] = userValidatedData.Data.Limit
	resp["offset"] = userValidatedData.Data.Offset
	return resp
}

func UserListRequestValidating(requestData map[string]any) UserListValidatedData {
	var response UserListValidatedData
	fmt.Println(requestData)
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
