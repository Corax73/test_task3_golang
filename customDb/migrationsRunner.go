package customDb

import (
	"checklist/customLog"
	"checklist/utils"
	"database/sql"
	"os"
	"strings"
)

func LoadSQLFile(fileName string) string {
	var resp string
	if fileName != "" {
		fileName = utils.ConcatSlice([]string{"./migrations/", fileName})
		_, err := os.Stat(fileName)
		if !os.IsNotExist(err) {
			file, err := (os.ReadFile(fileName))
			if err != nil {
				customLog.Logging(err)
			} else {
				resp = string(file)
			}
		} else {
			customLog.Logging(err)
		}
	}
	return resp
}

func RunTableMigration(db *sql.DB, tableName string) bool {
	var resp bool
	query := LoadSQLFile(utils.ConcatSlice([]string{tableName, "_up.sql"}))
	if query != "" {
		tx, err := db.Begin()
		if err != nil {
			customLog.Logging(err)
		} else {
			check := true
			for _, q := range strings.Split(string(query), ";") {
				q := strings.TrimSpace(q)
				if q == "" {
					continue
				}
				if _, err := tx.Exec(q); err != nil {
					customLog.Logging(err)
					tx.Rollback()
					check = false
					break
				}
			}
			if check {
				resp = true
				tx.Commit()
			}
		}
	}
	return resp
}
