package main

import (
	"checklist/customDb"
	"checklist/customLog"
)

func main() {
	customLog.LogInit("./logs/app.log")
	db := customDb.GetConnect()
	defer customDb.CloseConnect(db)
	customDb.RunTableMigration(db, "roles")
	customDb.RunTableMigration(db, "users")
}
