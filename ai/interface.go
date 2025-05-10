package ai

import "github.com/luizhenriquees/go-http-rpa/entity"

type ExamAssistant interface {
	GetAnswerIndex(question entity.Question) (int, error)
}
