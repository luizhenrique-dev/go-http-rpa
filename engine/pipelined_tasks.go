package engine

import (
	"errors"
	"fmt"
)

var (
	ErrNoTasksFound = errors.New("no tasks found")
)

type PipelinedTasks struct {
	tasks  []Task
	name   string
	Logger Logger
}

func NewPipelinedTasks(name string, tasks []Task) *PipelinedTasks {
	return &PipelinedTasks{
		name:   name,
		tasks:  tasks,
		Logger: &DefaultLogger{prefix: fmt.Sprintf("Pipelined Task - %s", name)},
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
	p.Logger.Info("Executing pipelined tasks...")
	for _, task := range p.tasks {
		if err := task.Execute(); err != nil {
			return fmt.Errorf("HTTP request failed in pipeline task '%s': %w", task.Name(), err)
		}
	}
	p.Logger.Info("Finished pipelined tasks")
	return nil
}

// AddTask adds a task to the pipeline
func (p *PipelinedTasks) AddTask(task Task) *PipelinedTasks {
	p.tasks = append(p.tasks, task)
	return p
}
