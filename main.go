package main

import (
	"checklist/customDb"
	"checklist/customLog"
	"checklist/router"
	"checklist/utils"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/savioxavier/termlink"
)

func main() {
	customLog.LogInit("./logs/app.log")
	db := customDb.GetConnect()
	defer customDb.CloseConnect(db)
	customDb.RunTableMigration(db, "roles")
	customDb.RunTableMigration(db, "users")
	customDb.RunTableMigration(db, "checklists")
	customDb.RunTableMigration(db, "checklist_items")
	router := (*&router.Router{}).Init()
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	defer close(errChan)
	mainPort := ":8080"
	envData := utils.GetConfFromEnvFile()
	if val, ok := envData["MAIN_PORT"]; ok {
		mainPort = utils.ConcatSlice([]string{":", val})
	}
	defer wg.Wait()
	wg.Add(1)
	go func(errChan chan<- error, handler http.Handler) {
		errChan <- http.ListenAndServe(mainPort, handler)
		defer wg.Done()
	}(errChan, router)

	go func() {
		check := true
		var invitationPrinted bool
		for check {
			if len(errChan) > 0 {
				fmt.Println(<-errChan)
				check = false
			} else {
				if !invitationPrinted {
					fmt.Println(strings.Join([]string{"started ",
						termlink.Link(
							utils.ConcatSlice([]string{"http://localhost", mainPort}),
							utils.ConcatSlice([]string{"http://localhost", mainPort}),
						)},
						" ",
					))
					invitationPrinted = true
				}
			}
		}
	}()
}
