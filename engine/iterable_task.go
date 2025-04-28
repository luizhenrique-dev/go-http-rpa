package engine

import (
	"errors"
	"fmt"
)

// CurrentElement represents a key or identifier for the current element in a collection or operation.
// CurrentIndex represents a key or identifier for the current index in a collection or operation.
const (
	CurrentElement = "currentElement"
	CurrentIndex   = "currentIndex"
)

// IterableTask represents a task that iterates over a collection of elements performing operations.
// It is a generic type parameterized by T, supporting various data types.
// The task leverages an inner Task implementation, using provided parameters to iterate.
// Elements are identified using a key and are fetched dynamically from the given parameters.
// Current iteration context, such as element and index, is stored in parameters for each iteration.
type IterableTask[T any] struct {
	task        Task
	elementsKey string
	params      Parameters
	name        string
	Logger      Logger
}

// NewIterableTask creates a new instance of IterableTask with a specified name, elementsKey, underlying task, and parameters.
// It enables iteration over elements of type T, processing each using the provided task.
func NewIterableTask[T any](name string, elementsKey string, task Task, params Parameters) *IterableTask[T] {
	return &IterableTask[T]{
		task,
		elementsKey,
		params,
		name,
		&DefaultLogger{prefix: fmt.Sprintf("Iterable Task - %s", name)},
	}
}

// Execute iterates over a list of elements, setting the current element and index, and executes the wrapped task for each item. Returns an error if the task execution fails.
func (t *IterableTask[T]) Execute() error {
	rawElements := t.params.Get(t.elementsKey)
	elements, ok := rawElements.([]T)
	if !ok {
		return fmt.Errorf("expected '%s' to be of type []T, but it was %T", t.elementsKey, rawElements)
	}
	t.Logger.Info("Initiating iterable task")
	for index, element := range elements {
		t.params.Put(t.elementsKey+"_"+CurrentElement, element)
		t.params.Put(t.elementsKey+"_"+CurrentIndex, index)
		t.Logger.Info("Executing task for element %s with index %d", element, index)
		if err := t.task.Execute(); err != nil {
			return fmt.Errorf("HTTP request failed in iterable task: %w", err)
		}
	}
	t.Logger.Info("Finished iterable task")
	return nil
}

// Validate checks if the IterableTask has a valid elementsKey and ensures the required parameter exists; returns an error if validation fails.
func (t *IterableTask[T]) Validate() error {
	if t.elementsKey == "" {
		return errors.New("missing element key")
	}
	if t.params.Get(t.elementsKey) == nil {
		return errors.New("missing required parameter: " + t.elementsKey)
	}
	return nil
}

// Name returns the name of the IterableTask. It is typically used for identification or debugging purposes.
func (t *IterableTask[T]) Name() string {
	return t.name
}
