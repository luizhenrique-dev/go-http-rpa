package rpaquiz

import (
	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

type TaskProcessQuiz struct {
	*engine.PipelinedTasks
}

func NewTaskProcessQuiz(headers httprequest.Headers, params engine.Parameters) *TaskProcessQuiz {
	pipeline := &TaskProcessQuiz{
		PipelinedTasks: engine.NewPipelinedTasks(
			"process_quiz",
			make([]engine.Task, 0, 2),
		),
	}

	pipeline.AddTask(NewTaskStartQuiz(headers, params))
	pipeline.AddTask(NewTaskAnswerQuestions(headers, params))
	return pipeline
}
