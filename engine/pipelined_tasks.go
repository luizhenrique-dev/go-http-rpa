package engine

import (
	"errors"
	"fmt"
)

var (
	ErrNoTasksFound = errors.New("no tasks found")
)

type PipelinedTasks struct {
	tasks []Task
	name  string
}

func NewPipelinedTasks(name string, tasks []Task) *PipelinedTasks {
	return &PipelinedTasks{
		name:  name,
		tasks: tasks,
	}
}

func (p *PipelinedTasks) Name() string {
	return p.name
}

func (p *PipelinedTasks) Validate() error {
	if len(p.tasks) == 0 {
		return ErrNoTasksFound
	}
	return nil
}

func (p *PipelinedTasks) Execute() error {
	for _, task := range p.tasks {
		if err := task.Execute(); err != nil {
			return fmt.Errorf("HTTP request failed in pipeline task '%s': %w", task.Name(), err)
		}
	}
	return nil
}

// AddTask adds a task to the pipeline
func (p *PipelinedTasks) AddTask(task Task) *PipelinedTasks {
	p.tasks = append(p.tasks, task)
	return p
}
