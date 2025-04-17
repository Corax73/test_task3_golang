package customStructs

type SimpleResponse struct {
	Success bool
	Message map[string]any
}

type ListResponse struct {
	Success bool
	Message []map[string]any
}
