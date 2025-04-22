package router

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/middleware"
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
	r.HandleFunc("/users/", router.getUsers).Methods("GET")

	r.HandleFunc("/roles/", router.createRole).Methods("POST")
	r.HandleFunc("/roles/", router.getRoles).Methods("GET")

	r.HandleFunc("/checklists/", router.createChecklist).Methods("POST")
	r.HandleFunc("/checklists/", router.getChecklists).Methods("GET")
	return &Router{r}
}

// initProcess returns a map of request parameters, causes console output on request.
func (router *Router) initProcess(w http.ResponseWriter, r *http.Request, getPost bool) map[string]any {
	resp := make(map[string]any)
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
	var response customStructs.SimpleResponse
	params := router.initProcess(w, r, true)
	response.Message = make(map[string]any, len(params))
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

// getUsers returns a list of entities, can use limit and offset parameters.
func (router *Router) getUsers(w http.ResponseWriter, r *http.Request) {
	login, passHash, ok := r.BasicAuth()
	if ok && middleware.BasicCheck(login, passHash) {
		var response customStructs.ListResponse
		params := router.initProcess(w, r, false)
		validatedData := validations.EntityListRequestValidating(params)
		userModel := (*&models.User{}).Init()
		if validatedData.Success {
			response.Message = userModel.GetList(validatedData.ToMap())
		} else {
			response.Message = userModel.GetList(make(map[string]string, 1))
		}
		if len(response.Message) > 1 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// createRole by post parameters creates an entity
func (router *Router) createRole(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	params := router.initProcess(w, r, true)
	response.Message = make(map[string]any, len(params))
	validatedData := validations.RoleCreateRequestValidating(params)
	if validatedData.Success {
		roleModel := (*&models.Role{}).Init()
		result := roleModel.Create(map[string]string{
			"id":         "",
			"title":      validatedData.Data.Title,
			"created_at": "",
		})
		if id, ok := result["id"]; !ok {
			response.Message["error"] = "Error.Try again"
		} else {
			response.Success = true
			response.Message["id"] = id
		}
	} else {
		response.Message = validatedData.ToMap()
	}

	json.NewEncoder(w).Encode(response)
}

// getRoles returns a list of entities, can use limit and offset parameters.
func (router *Router) getRoles(w http.ResponseWriter, r *http.Request) {
	var response customStructs.ListResponse
	params := router.initProcess(w, r, false)
	validatedData := validations.EntityListRequestValidating(params)
	roleModel := (*&models.Role{}).Init()
	if validatedData.Success {
		response.Message = roleModel.GetList(validatedData.ToMap())
	} else {
		response.Message = roleModel.GetList(make(map[string]string, 1))
	}
	if len(response.Message) > 1 {
		response.Success = true
	}
	json.NewEncoder(w).Encode(response)
}

// createChecklist by post parameters creates an entity
func (router *Router) createChecklist(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	params := router.initProcess(w, r, true)
	response.Message = make(map[string]any, len(params))
	validatedData := validations.ChecklistCreateRequestValidating(params)
	if validatedData.Success {
		checklistModel := (*&models.Checklist{}).Init()
		result := checklistModel.Create(map[string]string{
			"id":         "",
			"user_id":    validatedData.Data.UserId,
			"title":      validatedData.Data.Title,
			"created_at": "",
		})
		if id, ok := result["id"]; !ok {
			response.Message["error"] = "Error.Try again"
		} else {
			response.Success = true
			response.Message["id"] = id
		}
	} else {
		response.Message = validatedData.ToMap()
	}

	json.NewEncoder(w).Encode(response)
}

// getChecklists returns a list of entities, can use limit and offset parameters.
func (router *Router) getChecklists(w http.ResponseWriter, r *http.Request) {
	var response customStructs.ListResponse
	params := router.initProcess(w, r, false)
	validatedData := validations.EntityListRequestValidating(params)
	roleModel := (*&models.Role{}).Init()
	if validatedData.Success {
		response.Message = roleModel.GetList(validatedData.ToMap())
	} else {
		response.Message = roleModel.GetList(make(map[string]string, 1))
	}
	if len(response.Message) > 1 {
		response.Success = true
	}
	json.NewEncoder(w).Encode(response)
}
