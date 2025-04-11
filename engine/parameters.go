package engine

// Parameters holds the parameters for tasks
type Parameters map[string]any

func (params Parameters) Put(key string, value any) {
	params[key] = value
}

func (params Parameters) Get(key string) any {
	return params[key]
}
