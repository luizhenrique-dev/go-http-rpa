package rpaquiz

import (
	"github.com/luizhenriquees/go-http-rpa/engine"
	"github.com/luizhenriquees/go-http-rpa/entity"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const questionsKey = "questions"

type TaskAnswerQuestions struct {
	*engine.IterableTask[entity.Question]
}

func NewTaskAnswerQuestions(headers httprequest.Headers, params engine.Parameters) *TaskAnswerQuestions {
	return &TaskAnswerQuestions{
		IterableTask: engine.NewIterableTask[entity.Question](
			"answer_questions",
			questionsKey,
			NewTaskAnswerQuestion(headers, params),
			params,
		),
	}
}
