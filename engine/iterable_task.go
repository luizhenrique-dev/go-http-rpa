package engine

import (
	"errors"
	"fmt"
)

const (
	CurrentElement = "currentElement"
	CurrentIndex   = "currentIndex"
)

// IterableTask represents an iterable task, it fetches elements based on a key and iterates through them
type IterableTask[T any] struct {
	task        Task
	elementsKey string
	params      Parameters
	name        string
}

// NewIterableTask creates a new iterable HTTP task
func NewIterableTask[T any](name string, elementsKey string, task Task, params Parameters) *IterableTask[T] {
	return &IterableTask[T]{
		task,
		elementsKey,
		params,
		name,
	}
}

// Execute performs the HTTP request
func (t *IterableTask[T]) Execute() error {
	rawElements := t.params.Get(t.elementsKey)
	elements, ok := rawElements.([]T)
	if !ok {
		return fmt.Errorf("expected '%s' to be of type []T, but it was %T", t.elementsKey, rawElements)
	}
	for index, element := range elements {
		t.params.Put(CurrentElement, element)
		t.params.Put(CurrentIndex, index)
		if err := t.task.Execute(); err != nil {
			return fmt.Errorf("HTTP request failed in iterable task: %w", err)
		}
	}
	return nil
}

func (t *IterableTask[T]) Validate() error {
	if t.elementsKey == "" {
		return errors.New("missing element key")
	}
	if t.params.Get(t.elementsKey) == nil {
		return errors.New("missing required parameter: " + t.elementsKey)
	}
	return nil
}

func (t *IterableTask[T]) Name() string {
	return t.name
}
