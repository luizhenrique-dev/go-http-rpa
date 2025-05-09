package entity

type CoursesList struct {
	Courses []Course `json:"courses"`
}

type Course struct {
	ID        int      `json:"id"`
	TaskCount int      `json:"tasks_count,omitempty"`
	TaskDone  int      `json:"tasks_done,omitempty"`
	Status    string   `json:"status"`
	Modules   []Module `json:"modules"`
}

type Module struct {
	ID                int    `json:"id"`
	TaskCount         int    `json:"tasks_count,omitempty"`
	FinishedTaskCount int    `json:"finished_tasks_count,omitempty"`
	Tasks             []Task `json:"tasks"`
}

func (m *Module) IsFinished() bool {
	return m.TaskCount == m.FinishedTaskCount
}

type Task struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	CourseID       int        `json:"course_id"`
	ModuleID       int        `json:"module_id"`
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	QuestionsCount int        `json:"questions_count,omitempty"`
	Questions      []Question `json:"questions"`
}
