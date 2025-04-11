package engine

// Task represents a single operation within a rpa
type Task interface {
	Execute() error
	Validate() error
	Name() string
}
