package rpaquiz

import (
	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const quizIdsKey = "quizIds"

type TaskProcessQuizes struct {
	*engine.IterableTask[string]
}

func NewTaskProcessQuizes(headers httprequest.Headers, params engine.Parameters) *TaskProcessQuizes {
	return &TaskProcessQuizes{
		IterableTask: engine.NewIterableTask[string](
			"process_quizes",
			quizIdsKey,
			NewTaskProcessQuiz(headers, params),
			params,
		),
	}
}
