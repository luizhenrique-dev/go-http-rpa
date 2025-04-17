package rpaquiz

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

// TaskAnswerQuiz is a task to answer a quiz question
type TaskAnswerQuiz struct {
	*engine.HTTPTask
}

func NewTaskAnswerQuiz(headers httprequest.Headers, params engine.Parameters) *TaskAnswerQuiz {
	return &TaskAnswerQuiz{
		HTTPTask: engine.NewHTTPTask(
			"answer_quiz",
			httprequest.POST,
			headers,
			params,
		),
	}
}

func (t *TaskAnswerQuiz) Execute() error {
	baseURL := t.Params["baseUrl"].(string)
	quizID := t.Params[engine.CurrentElement].(string)
	quizData := t.Params["quizData"].(*QuizData)
	answerURL := baseURL + quizPath + quizID + "/answer"

	t.Logger.Info("Answering questions at: %s", answerURL)

	for index, question := range quizData.Questions {
		payload := t.createAnswerPayload(len(question.Options), index)

		t.Logger.Info("Question: %s", question.Question)
		t.Logger.Info("Number of possible answers: %d", len(question.Options))
		t.Logger.Info("Selected answer: %s", payload)
		resp, err := doPostRequest(answerURL, t.Headers, []byte(payload))
		if err != nil {
			return fmt.Errorf("failed to submit answer for question %d: %w", index, err)
		}

		t.Logger.Info("Question %d result - Correct: %d, Answered: %d",
			index, resp.Questions[index].Correct, resp.Questions[index].Answered)
		t.Logger.Info("=====================================================")
		time.Sleep(DefaultWaitTime)
	}

	return nil
}

func (t *TaskAnswerQuiz) createAnswerPayload(optionsCount int, questionIndex int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	answerIndex := rng.Intn(optionsCount)
	return fmt.Sprintf(`{"answer":%d,"question_index":%d}`, answerIndex, questionIndex)
}
