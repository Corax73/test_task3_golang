package router

import (
	"checklist/customLog"
	"checklist/customStructs"
	"checklist/middleware"
	"checklist/models"
	"checklist/utils"
	"checklist/validations"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	*mux.Router
}

func (router *Router) Init() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/", router.createUser).Methods("POST")
	r.HandleFunc("/users/", router.getUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", router.updateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", router.deleteUser).Methods("DELETE")

	r.HandleFunc("/roles/", router.createRole).Methods("POST")
	r.HandleFunc("/roles/", router.getRoles).Methods("GET")
	r.HandleFunc("/roles/{id:[0-9]+}", router.deleteRole).Methods("DELETE")

	r.HandleFunc("/checklists/", router.createChecklist).Methods("POST")
	r.HandleFunc("/checklists/", router.getChecklists).Methods("GET")
	r.HandleFunc("/checklists/{id:[0-9]+}", router.updateChecklist).Methods("PUT")
	r.HandleFunc("/checklists/{id:[0-9]+}", router.deleteChecklist).Methods("DELETE")

	r.HandleFunc("/checklists/items/", router.createChecklistItems).Methods("POST")
	r.HandleFunc("/checklists/{id:[0-9]+}/items/", router.getChecklistsItems).Methods("GET")
	r.HandleFunc("/checklists/items/{id:[0-9]+}", router.updateChecklistItem).Methods("PUT")
	r.HandleFunc("/checklists/items/{id:[0-9]+}", router.deleteChecklistItem).Methods("DELETE")
	return &Router{r}
}

// initProcess returns a map of request parameters, causes console output on request.
func (router *Router) initProcess(w http.ResponseWriter, r *http.Request, getPost bool, entityName, action string) customStructs.Request {
	login, passHash, ok := r.BasicAuth()
	var resp customStructs.Request
	resp.Params = make(map[string]any)
	if ok {
		passwordCheck, userData := middleware.BasicCheck(login, passHash)
		if passwordCheck && len(userData) > 0 && entityName != "" && action != "" {
			resp.Auth = middleware.UserCan(userData, entityName, action)
		} else {
			if passwordCheck {
				resp.Auth = true
			}
		}
		if resp.Auth {
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
				err := json.NewDecoder(r.Body).Decode(&resp.Params)
				if err != nil {
					customLog.Logging(err)
				}
			}
			for k, v := range vars {
				resp.Params[k] = v
			}
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

// setUnauthorized
func (router *Router) setUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// createUser by post parameters creates an entity
func (router *Router) createUser(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "users", "create")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.UserCreateRequestValidating(request)
		if validatedData.Success {
			passHash, err := utils.HashPassword(validatedData.Data.Password)
			if err == nil {
				userModel := (*&models.User{}).Init()
				result := userModel.Create(map[string]string{
					"id":                  "",
					"login":               validatedData.Data.Login,
					"role_id":             validatedData.Data.RoleId,
					"email":               validatedData.Data.Email,
					"password":            passHash,
					"checklists_quantity": validatedData.Data.ChecklistsQuantity,
					"created_at":          "",
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
	} else {
		router.setUnauthorized(w)
	}
}

// updateUser updates entity.
func (router *Router) updateUser(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "users", "update")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.UserUpdateRequestValidating(request)
		if validatedData.Success {
			userModel := (*&models.User{}).Init()
			result := userModel.Update(validatedData.ToMap(), validatedData.Data.Id)
			if id, ok := result["id"]; !ok {
				response.Message["error"] = "Error.Try again"
			} else {
				response.Success = true
				response.Message["id"] = id
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}

// getUsers returns a list of entities, can use limit and offset parameters.
func (router *Router) getUsers(w http.ResponseWriter, r *http.Request) {
	var response customStructs.ListResponse
	request := router.initProcess(w, r, false, "users", "read")
	if request.Auth {
		validatedData := validations.EntityListRequestValidating(request)
		userModel := (*&models.User{}).Init()
		if validatedData.Success {
			client := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			})
			ctx := context.Background()
			val, err := client.Get(ctx, validatedData.GetAsKey()).Result()
			if err != nil {
				customLog.Logging(err)
				response.Message = userModel.GetList(validatedData.ToMap())
				err := client.Set(ctx, validatedData.GetAsKey(), response.ToString(), 0).Err()
				if err != nil {
					customLog.Logging(err)
				}
			} else {
				if val != "" {
					jsonMap := make([]map[string]any, 100)
					err := json.Unmarshal([]byte(val), &jsonMap)
					if err != nil {
						customLog.Logging(err)
					} else {
						response.Message = jsonMap
					}
				} else {
					response.Message = userModel.GetList(validatedData.ToMap())
					err := client.Set(ctx, validatedData.GetAsKey(), response.ToString(), 0).Err()
					if err != nil {
						customLog.Logging(err)
					}
				}
			}
		} else {
			response.Message = userModel.GetList(make(map[string]string, 1))
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// deleteUser deletes an entity using the parameter `id`.
func (router *Router) deleteUser(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "users", "delete")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.EntityDeleteRequestValidating(request, "users")
		if validatedData.Success {
			userModel := (*&models.User{}).Init()
			userIdInt, _ := strconv.Atoi(validatedData.Data.Id)
			response.Message = userModel.Delete(userIdInt)
		} else {
			response.Message["error"] = "Error.Try again"
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// createRole by post parameters creates an entity
func (router *Router) createRole(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "roles", "create")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.RoleCreateRequestValidating(request)
		if validatedData.Success {
			roleModel := (*&models.Role{}).Init()
			result := roleModel.Create(map[string]string{
				"id":         "",
				"title":      validatedData.Data.Title,
				"abilities":  validatedData.Data.Abilities,
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
	} else {
		router.setUnauthorized(w)
	}
}

// getRoles returns a list of entities, can use limit and offset parameters.
func (router *Router) getRoles(w http.ResponseWriter, r *http.Request) {
	var response customStructs.ListResponse
	request := router.initProcess(w, r, false, "roles", "read")
	if request.Auth {
		validatedData := validations.EntityListRequestValidating(request)
		roleModel := (*&models.Role{}).Init()
		if validatedData.Success {
			response.Message = roleModel.GetList(validatedData.ToMap())
		} else {
			response.Message = roleModel.GetList(make(map[string]string, 1))
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// deleteRole deletes an entity using the parameter `id`.
func (router *Router) deleteRole(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "roles", "delete")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.EntityDeleteRequestValidating(request, "roles")
		if validatedData.Success {
			roleModel := (*&models.Role{}).Init()
			roleIdInt, _ := strconv.Atoi(validatedData.Data.Id)
			response.Message = roleModel.Delete(roleIdInt)
		} else {
			response.Message["error"] = "Error.Try again"
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// createChecklist by post parameters creates an entity
func (router *Router) createChecklist(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "checklists", "create")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.ChecklistCreateRequestValidating(request)
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
	} else {
		router.setUnauthorized(w)
	}
}

// getChecklists returns a list of entities, can use limit and offset parameters.
func (router *Router) getChecklists(w http.ResponseWriter, r *http.Request) {
	var response customStructs.ListResponse
	request := router.initProcess(w, r, false, "checklists", "read")
	if request.Auth {
		validatedData := validations.EntityListRequestValidating(request)
		checklistModel := (*&models.Checklist{}).Init()
		if validatedData.Success {
			response.Message = checklistModel.GetList(validatedData.ToMap())
		} else {
			response.Message = checklistModel.GetList(make(map[string]string, 1))
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// updateChecklist updates entity.
func (router *Router) updateChecklist(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "checklists", "update")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.ChecklistUpdateRequestValidating(request)
		if validatedData.Success {
			checklistModel := (*&models.Checklist{}).Init()
			result := checklistModel.Update(validatedData.ToMap(), validatedData.Data.Id)
			if id, ok := result["id"]; !ok {
				response.Message["error"] = "Error.Try again"
			} else {
				response.Success = true
				response.Message["id"] = id
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}

// deleteChecklist deletes an entity using the parameter `id`.
func (router *Router) deleteChecklist(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "checklists", "delete")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.EntityDeleteRequestValidating(request, "checklists")
		if validatedData.Success {
			checklistModel := (*&models.Checklist{}).Init()
			checklistIdInt, _ := strconv.Atoi(validatedData.Data.Id)
			response.Message = checklistModel.Delete(checklistIdInt)
		} else {
			response.Message["error"] = "Error.Try again"
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// createChecklistItems by post parameters creates an entity
func (router *Router) createChecklistItems(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "checklist_items", "create")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.ChecklistItemCreateRequestValidating(request)
		if validatedData.Success {
			checklistModel := (*&models.ChecklistItem{}).Init()
			result := checklistModel.Create(map[string]string{
				"id":           "",
				"checklist_id": validatedData.Data.ChecklistId,
				"is_completed": validatedData.Data.IsCompleted,
				"description":  validatedData.Data.Description,
				"created_at":   "",
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
	} else {
		router.setUnauthorized(w)
	}
}

// getChecklistsItems returns a list of entities, can use limit and offset parameters.
func (router *Router) getChecklistsItems(w http.ResponseWriter, r *http.Request) {
	var response customStructs.ListResponse
	request := router.initProcess(w, r, false, "checklist_items", "read")
	if request.Auth {
		validatedData := validations.EntityListRequestValidating(request)
		checklistModel := (*&models.ChecklistItem{}).Init()
		if validatedData.Success {
			if validatedData.Data.Id != "" {
				validatedData.Data.FilterBy = "checklist_id"
				validatedData.Data.FilterVal = validatedData.Data.Id
			}
			response.Message = checklistModel.GetList(validatedData.ToMap())
		} else {
			response.Message = checklistModel.GetList(make(map[string]string, 1))
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// deleteChecklistItem deletes an entity using the parameter `id`.
func (router *Router) deleteChecklistItem(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "checklist_items", "delete")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.EntityDeleteRequestValidating(request, "checklist_items")
		if validatedData.Success {
			checklistItemModel := (*&models.ChecklistItem{}).Init()
			checklistItemIdInt, _ := strconv.Atoi(validatedData.Data.Id)
			response.Message = checklistItemModel.Delete(checklistItemIdInt)
		} else {
			response.Message["error"] = "Error.Try again"
		}
		if len(response.Message) > 0 {
			response.Success = true
		}
		json.NewEncoder(w).Encode(response)
	} else {
		router.setUnauthorized(w)
	}
}

// updateChecklistItem updates entity.
func (router *Router) updateChecklistItem(w http.ResponseWriter, r *http.Request) {
	var response customStructs.SimpleResponse
	request := router.initProcess(w, r, true, "checklist_items", "update")
	if request.Auth {
		response.Message = make(map[string]any, len(request.Params))
		validatedData := validations.ChecklistItemUpdateRequestValidating(request)
		if validatedData.Success {
			checklistItemModel := (*&models.ChecklistItem{}).Init()
			result := checklistItemModel.Update(validatedData.ToMap(), validatedData.Data.Id)
			if id, ok := result["id"]; !ok {
				response.Message["error"] = "Error.Try again"
			} else {
				response.Success = true
				response.Message["id"] = id
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}
