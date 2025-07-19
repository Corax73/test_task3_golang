package customStructs

type Request struct {
	Auth   bool
	Params map[string]any
	Filters []map[string]string
}
