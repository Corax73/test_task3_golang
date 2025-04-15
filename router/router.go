package router

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/models"
	"checklist/utils"
	"checklist/validations"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func (router *Router) Init() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/", router.createUser).Methods("POST")
	return &Router{r}
}

// initProcess returns a map of request parameters, causes console output on request.
func (router *Router) initProcess(w http.ResponseWriter, r *http.Request, getPost bool) map[string]interface{} {
	var resp map[string]interface{}
	var vars map[string]string
	w.Header().Set("Content-Type", "application/json")
	if router.checkEnv() {
		router.consoleOutput(r)
	}
	vars = mux.Vars(r)
	sort := r.URL.Query().Get("sort")
	if sort != "" {
		var order string
		splits := strings.Split(sort, "--")
		if len(splits) > 1 {
			requestField, requestOrder := splits[0], splits[1]
			if requestOrder != "desc" && requestOrder != "asc" {
				order = "desc"
			} else {
				order = requestOrder
			}
			vars["order"] = order
			vars["orderBy"] = requestField
		}
	}
	filter := r.URL.Query().Get("filter")
	if filter != "" {
		splits := strings.Split(filter, "--")
		if len(splits) > 1 {
			requestField, requestValue := splits[0], splits[1]
			if requestValue != "" {
				vars["filterBy"] = requestField
				vars["filterVal"] = requestValue
			}
		}
	}
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		vars["limit"] = limit
	}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		vars["offset"] = offset
	}
	if getPost {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			customLog.Logging(err)
		}
	}
	for k, v := range vars {
		resp[k] = v
	}
	return resp
}

// checkEnv looks for a key `CONSOLE_OUT` in the .env file and returns true if its value is true.
func (router *Router) checkEnv() bool {
	var resp bool
	envData := utils.GetConfFromEnvFile()
	if val, ok := envData["CONSOLE_OUT"]; ok && val == "true" {
		resp = true
	}
	return resp
}

// consoleOutput displays the time, route and request method to the console.
func (router *Router) consoleOutput(r *http.Request) {
	fmt.Println(strings.Join([]string{time.Now().Format(time.RFC3339), r.Method, r.RequestURI, r.UserAgent()}, " "))
}

// createUser by post parameters creates an entity
func (router *Router) createUser(w http.ResponseWriter, r *http.Request) {
	var response customStructs.Response
	params := router.initProcess(w, r, true)
	response.Message = make(map[string]interface{}, len(params))
	validatedData := validations.UserCreateRequestValidating(params)
	if validatedData.Success {
		passHash, err := utils.HashPassword(validatedData.Data.Password)
		if err == nil {
			userModel := (*&models.User{}).Init()
			result := userModel.Create(map[string]string{
				"id":         "",
				"login":      validatedData.Data.Login,
				"role_id":    validatedData.Data.RoleId,
				"email":      validatedData.Data.Email,
				"password":   passHash,
				"created_at": "",
			})
			if id, ok := result["id"]; !ok {
				response.Message["error"] = "Error.Try again"
			} else {
				response.Success = true
				response.Message["id"] = id
			}
		} else {
			customLog.Logging(err)
			response.Message["password"] = "Please try again, if the error persists, contact the administrator"
		}
	} else {
		response.Message = validatedData.ToMap()
	}

	json.NewEncoder(w).Encode(response)
}
