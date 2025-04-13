package router

import (
	"checklist/customLog"
	"checklist/models"
	"checklist/utils"
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

// @var response stub for answers on routes.
var response map[string]interface{}

func (router *Router) Init() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/", router.createUser).Methods("POST")
	return &Router{r}
}

// initProcess returns a map of request parameters, causes console output on request.
func (router *Router) initProcess(w http.ResponseWriter, r *http.Request, getPost bool) map[string]string {
	var resp map[string]string
	w.Header().Set("Content-Type", "application/json")
	if router.checkEnv() {
		router.consoleOutput(r)
	}
	resp = mux.Vars(r)
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
			resp["order"] = order
			resp["orderBy"] = requestField
		}
	}
	filter := r.URL.Query().Get("filter")
	if filter != "" {
		splits := strings.Split(filter, "--")
		if len(splits) > 1 {
			requestField, requestValue := splits[0], splits[1]
			if requestValue != "" {
				resp["filterBy"] = requestField
				resp["filterVal"] = requestValue
			}
		}
	}
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		resp["limit"] = limit
	}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		resp["offset"] = offset
	}
	if getPost {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			customLog.Logging(err)
		}
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

// createUser by post parameters creates an entity, you, if the model implements the interface,
// then a request is made to enrich the entity data.
func (router *Router) createUser(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{}
	params := router.initProcess(w, r, true)
	if login, ok := params["login"]; ok && login != "" {
		userModel := (*&models.User{}).Init()
		if email, ok := params["email"]; ok && email != "" {
			if password, ok := params["password"]; ok && password != "" {
				result := userModel.Create(map[string]string{
					"id":         "",
					"login":      login,
					"role_id":    "1", // @todo
					"email":      email,
					"password":   password,
					"created_at": "",
				})
				if id, ok := result["id"]; !ok {
					response["data"] = "Error.Try again"
				} else {
					response["id"] = id
				}
			}
		} else {
			response["error"] = "Check parameters"
		}
	}
	json.NewEncoder(w).Encode(response)
}
