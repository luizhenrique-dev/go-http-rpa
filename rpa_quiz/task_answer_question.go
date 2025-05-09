package rpaquiz

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	"github.com/luizhenriquees/go-http-rpa/entity"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

// TaskAnswerQuestion is a task to answer a question from quiz
type TaskAnswerQuestion struct {
	*engine.HTTPTask
}

func NewTaskAnswerQuestion(headers httprequest.Headers, params engine.Parameters) *TaskAnswerQuestion {
	task := &TaskAnswerQuestion{}
	httpTask := engine.NewHTTPTask(
		"answer_question",
		httprequest.POST,
		headers,
		params,
		engine.WithPreRequestFunc(task.preRequest),
		engine.WithPostExtractFunc(task.postExtract),
	)
	task.HTTPTask = httpTask
	return task
}

func (t *TaskAnswerQuestion) preRequest() error {
	answerURL := t.getAnswerURL()
	t.Logger.Info("pre-request built URL: %s", answerURL)
	t.URL = answerURL

	var currentQuestion entity.Question
	if val, ok := t.Params.Get(questionsKey + "_" + engine.CurrentElement).(entity.Question); ok {
		currentQuestion = val
	}
	index := 0
	if val, ok := t.Params.Get(questionsKey + "_" + engine.CurrentIndex).(int); ok {
		index = val
	}
	payload := t.createAnswerPayload(len(currentQuestion.Options), index)
	t.RequestBody = []byte(payload)

	t.Logger.Info("Question %d: %s", index, currentQuestion.Question)
	t.Logger.Info("Number of possible answers: %d", len(currentQuestion.Options))
	t.Logger.Info("Selected answer: %s", payload)
	return nil
}

func (t *TaskAnswerQuestion) getAnswerURL() string {
	baseURL := ""
	if val, ok := t.Params.Get(engine.ParamBaseURL).(string); ok {
		baseURL = val
	}
	quizID := ""
	if val, ok := t.Params.Get(quizIdsKey + "_" + engine.CurrentElement).(string); ok {
		quizID = val
	}
	answerURL := baseURL + quizPath + quizID + "/answer"
	return answerURL
}

func (t *TaskAnswerQuestion) createAnswerPayload(optionsCount int, questionIndex int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	answerIndex := rng.Intn(optionsCount)
	return fmt.Sprintf(`{"answer":%d,"question_index":%d}`, answerIndex, questionIndex)
}

func (t *TaskAnswerQuestion) postExtract(resp *http.Response, _ *engine.HTTPTask) error {
	t.Logger.Info("PostExtract used.")
	defer resp.Body.Close()
	var responseData entity.QuizData
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response not OK: %s", resp.Status)
	}

	index := 0
	if val, ok := t.Params.Get(questionsKey + "_" + engine.CurrentIndex).(int); ok {
		index = val
	}
	t.Logger.Info("Question %d result - Correct: %d, Answered: %d",
		index, responseData.Questions[index].Correct, responseData.Questions[index].Answered)
	t.Logger.Info("=====================================================")
	return nil
}
