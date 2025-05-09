package models

import (
	"checklist/customDb"
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/utils"
	"slices"
	"strings"
)

type Model struct {
	table         string
	Fields        map[string]string
	GuardedFields []string
}

func (model *Model) Table() string {
	return model.table
}

func (model *Model) SetTable(tableTitle string) {
	model.table = tableTitle
}

func (model *Model) Create(fields map[string]string) map[string]string {
	response := map[string]string{}
	if utils.CompareMapsByStringKeys(model.Fields, fields) {
		model.Fields = fields
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		response = model.Save()
	}
	return response
}

func (model *Model) Save() map[string]string {
	response := map[string]string{}
	if len(model.Fields) > 0 {
		strSlice := make([]string, 5+((len(model.Fields)-1)*2))
		strSlice = append(strSlice, "INSERT INTO ")
		strSlice = append(strSlice, model.Table())
		strSlice = append(strSlice, " (")
		fields := utils.GetMapKeysWithValue(model.Fields)
		index := utils.GetIndexByStrValue(fields, "id")
		if index != -1 {
			fields = slices.Delete(fields, index, index+1)
		}
		strSlice = append(strSlice, strings.Trim(strings.Join(fields, ","), ","))
		strSlice = append(strSlice, ") VALUES (")
		values := utils.GetMapValues(model.Fields)
		valuesToDb := make([]string, len(values))
		for _, val := range fields {
			if _, ok := model.Fields[val]; ok {
				value := model.Fields[val]
				if strings.Contains(value, "'") {
					value = strings.Replace(value, "'", "''", -1)
				}
				valuesToDb = append(valuesToDb, utils.ConcatSlice([]string{"'", value, "'"}))
			}
		}
		strSlice = append(strSlice, strings.Trim(strings.Join(valuesToDb, ","), ","))
		strSlice = append(strSlice, ") RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		var id string
		err := db.QueryRow(queryStr).Scan(&id)
		if err != nil {
			customLog.Logging(err)
		} else {
			response = map[string]string{"id": id}
		}
	}
	return response
}

func (model *Model) GetList(params map[string]string) []map[string]any {
	var resp []map[string]any
	db := customDb.GetConnect()
	defer customDb.CloseConnect(db)
	modelFields := utils.GetMapKeys(model.Fields)
	selectedStr := ""
	for _, val := range modelFields {
		if !slices.Contains(model.GuardedFields, val) {
			selectedStr = utils.ConcatSlice([]string{
				selectedStr,
				val,
				", ",
			})
		}
	}
	selectedStr = strings.Trim(selectedStr, ", ")
	queryStr := utils.ConcatSlice([]string{
		"SELECT ",
		selectedStr,
		" FROM ",
		model.Table(),
	})
	if len(params) > 0 {
		if filterBy, ok := params["filterBy"]; ok && filterBy != "" {
			queryStr = utils.ConcatSlice([]string{
				queryStr,
				" WHERE ",
				params["filterBy"],
				" = '",
				params["filterVal"],
				"'",
			})
		}
		if order, ok := params["order"]; ok && order != "" {
			queryStr = utils.ConcatSlice([]string{
				queryStr,
				" ORDER BY ",
				params["orderBy"],
				" ",
				params["order"],
			})
		}
		if limit, ok := params["limit"]; ok && limit != "" {
			queryStr = utils.ConcatSlice([]string{
				queryStr,
				" LIMIT ",
				params["limit"],
			})
		}
		if offset, ok := params["offset"]; ok && offset != "" {
			queryStr = utils.ConcatSlice([]string{
				queryStr,
				" OFFSET ",
				params["offset"],
			})
		}
	}
	queryStr = utils.ConcatSlice([]string{
		queryStr,
		" ;",
	})
	rows, err := db.Query(queryStr)
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = utils.SqlToMap(rows)
	}
	return resp
}

func (model *Model) GetOneById(id int) customStructs.SimpleResponse {
	var resp customStructs.SimpleResponse
	if id > 0 {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		queryStr := utils.ConcatSlice([]string{
			"SELECT * FROM ",
			model.Table(),
			" WHERE id=$1;",
		})
		rows, err := db.Query(queryStr, id)
		if err != nil {
			customLog.Logging(err)
		} else {
			if data := utils.SqlToMap(rows); len(data) > 0 {
				resp.Success = true
				resp.Message = data[0]
			}
		}
	}
	return resp
}

func (model *Model) GetOneByField(field, value, withRelation string) customStructs.SimpleResponse {
	var resp customStructs.SimpleResponse
	if field != "" && value != "" {
		fieldNames := utils.GetMapKeys(model.Fields)
		if slices.Contains(fieldNames, field) {
			db := customDb.GetConnect()
			defer customDb.CloseConnect(db)

			queryStr := utils.ConcatSlice([]string{
				"SELECT * FROM ",
				model.Table(),
			})
			if withRelation == "roles" {
				queryStr = utils.ConcatSlice([]string{
					queryStr,
					" JOIN roles",
					" ON roles.id = users.role_id",
				})

			}
			queryStr = utils.ConcatSlice([]string{
				queryStr,
				" WHERE ",
				field,
				" = $1;",
			})
			rows, err := db.Query(queryStr, value)
			if err != nil {
				customLog.Logging(err)
			} else {
				if data := utils.SqlToMap(rows); len(data) > 0 {
					resp.Success = true
					resp.Message = data[0]
				}
			}
		}
	}
	return resp
}

func (model *Model) Update(fields map[string]string, id string) map[string]string {
	response := map[string]string{}
	fields = utils.GetMapWithoutKeys(fields, []string{"id"})
	if utils.PresenceMapKeysInOtherMap(fields, model.Fields) {
		strSlice := make([]string, 7+((len(fields)-1)*2))
		strSlice = append(strSlice, "UPDATE ")
		strSlice = append(strSlice, model.Table())
		strSlice = append(strSlice, " SET ")
		if len(fields) > 1 {
			strSlice = append(strSlice, "(")
		}
		columns := utils.GetMapKeysWithValue(fields)
		index := utils.GetIndexByStrValue(columns, "id")
		if index != -1 {
			columns = slices.Delete(columns, index, index+1)
		}
		if len(fields) > 1 {
			strSlice = append(strSlice, strings.Trim(strings.Join(columns, ","), ","))
			strSlice = append(strSlice, ") = (")
		} else {
			strSlice = append(strSlice, columns[0])
			strSlice = append(strSlice, " = ")
		}
		valuesToDb := make([]string, len(columns))
		var i int
		for _, val := range columns {
			if _, ok := fields[val]; ok {
				value := fields[val]
				if strings.Contains(value, "'") {
					value = strings.Replace(value, "'", "''", -1)
				}
				valuesToDb[i] = utils.ConcatSlice([]string{"'", value, "'"})
			}
			i++
		}
		if len(fields) > 1 {
			strSlice = append(strSlice, strings.Trim(strings.Join(valuesToDb, ","), ","))
			strSlice = append(strSlice, ") ")
		} else {
			strSlice = append(strSlice, valuesToDb[0])
			strSlice = append(strSlice, " ")
		}
		strSlice = append(strSlice, "WHERE id = ")
		strSlice = append(strSlice, id)
		strSlice = append(strSlice, " RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		var id string
		err := db.QueryRow(queryStr).Scan(&id)
		if err != nil {
			customLog.Logging(err)
		} else {
			response = map[string]string{"id": id}
		}
	}
	return response
}

func (model *Model) Delete(id int) map[string]any {
	resp := map[string]any{"success": false, "error": "not found"}
	if id > 0 {
		db := customDb.GetConnect()
		defer customDb.CloseConnect(db)
		queryStr := utils.ConcatSlice([]string{
			"DELETE FROM ",
			model.Table(),
			" WHERE id=$1 RETURNING id;",
		})
		rows, err := db.Query(queryStr, id)
		if err != nil {
			customLog.Logging(err)
		} else {
			if data := utils.SqlToMap(rows); len(data) > 0 {
				resp = data[0]
			}
		}
	}
	return resp
}
