package rpa

import (
	"errors"
	"fmt"
)

// Parameters holds the parameters for tasks
type Parameters map[string]any

// Task represents a single operation within a rpa
type Task interface {
	Execute(params Parameters) error
	Validate(params Parameters) error
	Name() string
}

// Rpa represents a complete workflow consisting of multiple tasks
type Rpa struct {
	tasks  []Task
	params Parameters
	logger Logger
	name   string
}

// NewRpa creates a new rpa with the given name
func NewRpa(name string) *Rpa {
	return &Rpa{
		name:   name,
		tasks:  []Task{},
		logger: &DefaultLogger{},
		params: make(Parameters),
	}
}

// AddTask adds a task to the rpa
func (j *Rpa) AddTask(task Task) *Rpa {
	j.tasks = append(j.tasks, task)
	return j
}

// SetParams sets the parameters for the rpa
func (j *Rpa) SetParams(params Parameters) *Rpa {
	j.params = params
	return j
}

// AddParam adds a single parameter to the rpa
func (j *Rpa) AddParam(key string, value any) *Rpa {
	j.params[key] = value
	return j
}

// Execute runs all tasks in the rpa sequentially, validating each task before execution
func (j *Rpa) Execute() error {
	if j == nil {
		return errors.New("rpa is nil")
	}

	j.logger.Info("Starting rpa: %s", j.name)
	for _, task := range j.tasks {
		taskName := task.Name()

		j.logger.Info("Validating task: %s", taskName)
		if err := task.Validate(j.params); err != nil {
			j.logger.Error("Validation failed for task", err, "task", taskName)
			return fmt.Errorf("validation failed for task %s: %w", taskName, err)
		}

		j.logger.Info("Executing task: %s", taskName)
		if err := task.Execute(j.params); err != nil {
			j.logger.Error("Task execution failed", err, "task", taskName)
			return fmt.Errorf("execution failed for task %s: %w", taskName, err)
		}
	}

	j.logger.Info("Rpa completed successfully: %s", j.name)
	return nil
}

// BaseTask provides common functionality for tasks
type BaseTask struct {
	requiredParams []string
	name           string
}

// NewBaseTask creates a new base task with the given name
func NewBaseTask(name string, requiredParams ...string) BaseTask {
	return BaseTask{
		name:           name,
		requiredParams: requiredParams,
	}
}

// Name returns the name of the task
func (b BaseTask) Name() string {
	return b.name
}

// Validate checks if all required parameters exist
func (b BaseTask) Validate(params Parameters) error {
	for _, param := range b.requiredParams {
		if _, exists := params[param]; !exists {
			return errors.New("missing required parameter: " + param)
		}
	}
	return nil
}
