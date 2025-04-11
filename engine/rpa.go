package engine

import (
	"errors"
	"fmt"

	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	ParamBaseURL = "baseUrl"
)

// Rpa represents a complete workflow consisting of multiple tasks
type Rpa struct {
	tasks          []Task
	params         Parameters
	defaultHeaders httprequest.Headers
	logger         Logger
	name           string
	baseURL        string
}

// NewRpa creates a new rpa with the given name
func NewRpa(name, baseURL string, defaultHeaders httprequest.Headers) *Rpa {
	return &Rpa{
		name:           name,
		baseURL:        baseURL,
		tasks:          []Task{},
		logger:         &DefaultLogger{prefix: "RPA"},
		defaultHeaders: defaultHeaders,
		params:         make(Parameters),
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

func (j *Rpa) GetParams() Parameters {
	return j.params
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
		if err := task.Validate(); err != nil {
			j.logger.Error("Validation failed for task", err, "task", taskName)
			return fmt.Errorf("validation failed for task %s: %w", taskName, err)
		}

		j.logger.Info("Executing task: %s", taskName)
		if err := task.Execute(); err != nil {
			j.logger.Error("Task execution failed", err, "task", taskName)
			return fmt.Errorf("execution failed for task %s: %w", taskName, err)
		}
		j.logger.Info("Finishing task: %s", taskName)
	}

	j.logger.Info("Rpa completed successfully: %s", j.name)
	return nil
}
