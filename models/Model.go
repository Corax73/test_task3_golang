package models

import (
	"checklist/customDb"
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/utils"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Model struct {
	table              string
	Fields, FieldTypes map[string]string
	GuardedFields      []string
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
		valuesToDb := make([]any, len(values))
		valPlaceholdersSlice := make([]string, len(fields))
		var i int
		for _, val := range fields {
			if value, ok := model.Fields[val]; ok {
				if fieldType, ok := model.FieldTypes[val]; ok {
					if val == "created_at" {
						t, err := time.Parse(time.RFC3339, value)
						if err != nil {
							customLog.Logging(err)
						} else {
							valuesToDb[i] = t
						}
					} else {
						switch fieldType {
						case "bool":
							boolValue, err := strconv.ParseBool(value)
							if err != nil {
								customLog.Logging(err)
							} else {
								valuesToDb[i] = boolValue
							}
						case "int":
							intValue, err := strconv.Atoi(value)
							if err != nil {
								customLog.Logging(err)
							} else {
								valuesToDb[i] = intValue
							}
						default:
							valuesToDb[i] = value
						}
					}
					valPlaceholdersSlice = append(valPlaceholdersSlice, utils.ConcatSlice([]string{"$", strconv.Itoa(i + 1), ", "}))
				}
			}
			i++
		}
		strSlice = append(strSlice, strings.Trim(utils.ConcatSlice(valPlaceholdersSlice), ", "))
		strSlice = append(strSlice, ") RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		tx, err := db.Begin()
		defer customDb.CloseConnect(db)
		if err != nil {
			customLog.Logging(err)
		} else {
			defer tx.Rollback()
			stmt, err := tx.Prepare(queryStr)
			if err != nil {
				customLog.Logging(err)
			}
			defer stmt.Close()
			var id int
			err = stmt.QueryRow(valuesToDb...).Scan(&id)
			if err != nil {
				customLog.Logging(err)
			} else {
				err = tx.Commit()
				if err != nil {
					customLog.Logging(err)
				} else {
					response = map[string]string{"id": strconv.Itoa(id)}
				}
			}
		}
	}
	return response
}

func (model *Model) GetList(params map[string]string, additionalFilters []map[string]any) ([]map[string]any, int) {
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
	queryStrToTotal := utils.ConcatSlice([]string{
		"SELECT COUNT(id)",
		" FROM ",
		model.Table(),
	})
	valuesToDb := make([]any, len(additionalFilters))
	if len(params) > 0 {
		var hasMainFilter bool
		if filterBy, ok := params["filterBy"]; ok && filterBy != "" {
			hasMainFilter = true
			queryStr = utils.ConcatSlice([]string{
				queryStr,
				" WHERE ",
				params["filterBy"],
				" = '",
				params["filterVal"],
				"'",
			})
			queryStrToTotal = utils.ConcatSlice([]string{
				queryStrToTotal,
				" WHERE ",
				params["filterBy"],
				" = '",
				params["filterVal"],
				"'",
			})
		}
		if len(additionalFilters) > 0 {
			for i, filter := range additionalFilters {
				operator := " WHERE "
				if hasMainFilter {
					operator = " AND "
				}
				beging := ""
				ending := ""
				conditions := " = "
				if filter["conditions"] == "contains" {
					conditions = " LIKE "
					beging = "%"
					ending = "%"
				}
				if filter["conditions"] == "begin" {
					conditions = " LIKE "
					ending = "%"
				}
				if filter["conditions"] == "end" {
					conditions = " LIKE "
					beging = "%"
				}
				if filter["value"] != nil {
					switch filter["value"].(type) {
					case bool:
						valuesToDb[i] = filter["value"].(bool)
					case int:
						valuesToDb[i] = filter["value"].(int)
						conditions = " = "
					case int64:
						valuesToDb[i] = filter["value"].(int64)
						conditions = " = "
					case float64:
						valuesToDb[i] = filter["value"].(float64)
						conditions = " = "
					case string:
						valuesToDb[i] = utils.ConcatSlice([]string{beging, filter["value"].(string), ending})
					case time.Time:
						valuesToDb[i] = filter["value"].(time.Time)
					case []byte:
						valuesToDb[i] = string(filter["value"].([]byte))
					default:
						valuesToDb[i] = filter["value"]
					}
				}
				queryStr = utils.ConcatSlice([]string{
					queryStr,
					operator,
					filter["field"].(string),
					conditions,
					utils.ConcatSlice([]string{"$", strconv.Itoa(i + 1)}),
				})
				queryStrToTotal = utils.ConcatSlice([]string{
					queryStrToTotal,
					operator,
					filter["field"].(string),
					conditions,
					utils.ConcatSlice([]string{"$", strconv.Itoa(i + 1)}),
				})
				hasMainFilter = true
			}
			queryStr = strings.Trim(queryStr, ", ")
			queryStrToTotal = strings.Trim(queryStrToTotal, ", ")
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
	queryStrToTotal = utils.ConcatSlice([]string{
		queryStrToTotal,
		" ;",
	})
	var total int
	err := db.QueryRow(queryStrToTotal, valuesToDb...).Scan(&total)
	if err != nil {
		customLog.Logging(err)
	}
	queryStr = utils.ConcatSlice([]string{
		queryStr,
		" ;",
	})

	rows, err := db.Query(queryStr, valuesToDb...)
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = utils.SqlToMap(rows)
	}
	return resp, total
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
		strSlice := make([]string, 9+((len(fields)-1)*2))
		strSlice = append(strSlice, "UPDATE ")
		strSlice = append(strSlice, model.Table())
		strSlice = append(strSlice, " SET ")
		columns := utils.GetMapKeysWithValue(fields)
		if len(columns) > 1 {
			strSlice = append(strSlice, "(")
		}
		index := utils.GetIndexByStrValue(columns, "id")
		if index != -1 {
			columns = slices.Delete(columns, index, index+1)
		}
		if len(columns) > 1 {
			strSlice = append(strSlice, strings.Trim(strings.Join(columns, ","), ","))
			strSlice = append(strSlice, ") = (")
		} else {
			strSlice = append(strSlice, columns[0])
			strSlice = append(strSlice, " = ")
		}
		valuesToDb := make([]any, len(columns))
		var i int
		valPlaceholdersSlice := make([]string, len(columns))
		for _, val := range columns {
			if value, ok := fields[val]; ok {
				if fieldType, ok := model.FieldTypes[val]; ok {
					switch fieldType {
					case "bool":
						boolValue, err := strconv.ParseBool(value)
						if err != nil {
							customLog.Logging(err)
						} else {
							valuesToDb[i] = boolValue
						}
					case "int":
						intValue, err := strconv.Atoi(value)
						if err != nil {
							customLog.Logging(err)
						} else {
							valuesToDb[i] = intValue
						}
					default:
						value = strings.ReplaceAll(value, "'", "''")
						valuesToDb[i] = value
					}
					valPlaceholdersSlice = append(valPlaceholdersSlice, utils.ConcatSlice([]string{"$", strconv.Itoa(i + 1), ", "}))
					i++
				}
			}
		}
		strSlice = append(strSlice, strings.Trim(utils.ConcatSlice(valPlaceholdersSlice), ", "))
		if len(columns) > 1 {
			strSlice = append(strSlice, ") ")
		} else {
			strSlice = append(strSlice, " ")
		}

		strSlice = append(strSlice, utils.ConcatSlice([]string{"WHERE id = ", "$", strconv.Itoa(i + 1)}))
		strSlice = append(strSlice, " RETURNING id;")
		queryStr := utils.ConcatSlice(strSlice)
		db := customDb.GetConnect()
		tx, err := db.Begin()
		defer customDb.CloseConnect(db)
		if err != nil {
			customLog.Logging(err)
		} else {
			defer tx.Rollback()
			stmt, err := tx.Prepare(queryStr)
			if err != nil {
				customLog.Logging(err)
			}
			defer stmt.Close()
			valuesToDb = append(valuesToDb, id)
			row, err := stmt.Exec(valuesToDb...)
			if err != nil {
				customLog.Logging(err)
			} else {
				err = tx.Commit()
				if err != nil {
					customLog.Logging(err)
				} else {
					_, err := row.RowsAffected()
					if err != nil {
						customLog.Logging(err)
					} else {
						response = map[string]string{"id": id}
					}
				}
			}
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
