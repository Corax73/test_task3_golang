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
	table        string
	Fields       map[string]string
	FieldsStruct any
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

func (model *Model) GetList(params map[string]string) []map[string]interface{} {
	var resp []map[string]interface{}
	db := customDb.GetConnect()
	defer customDb.CloseConnect(db)
	queryStr := utils.ConcatSlice([]string{
		"SELECT * FROM ",
		model.Table(),
	})
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

func (model *Model) GetOneById(id int) customStructs.Response {
	var resp customStructs.Response
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
