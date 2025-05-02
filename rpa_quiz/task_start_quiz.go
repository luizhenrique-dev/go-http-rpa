package rpaquiz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

// TaskStartQuiz is a task to start a quiz
type TaskStartQuiz struct {
	*engine.HTTPTask
}

func NewTaskStartQuiz(headers httprequest.Headers, params engine.Parameters) *TaskStartQuiz {
	task := &TaskStartQuiz{}
	httpTask := engine.NewHTTPTask(
		"start_quiz",
		httprequest.POST,
		headers,
		params,
		engine.WithPreRequestFunc(task.preRequest),
		engine.WithPostExtractFunc(task.postExtract),
	)
	task.HTTPTask = httpTask
	return task
}

func (t *TaskStartQuiz) preRequest() error {
	baseURL := ""
	if val, ok := t.Params.Get(engine.ParamBaseURL).(string); ok {
		baseURL = val
	}
	quizID := ""
	if val, ok := t.Params.Get(quizIdsKey + "_" + engine.CurrentElement).(string); ok {
		quizID = val
	}
	startURL := baseURL + quizPath + quizID + "/start"
	t.Logger.Info("pre-request built URL: %s", startURL)
	t.URL = startURL
	return nil
}

func (t *TaskStartQuiz) postExtract(resp *http.Response, _ *engine.HTTPTask) error {
	t.Logger.Info("PostExtract used.")
	defer resp.Body.Close()
	var quizData QuizData
	if err := json.NewDecoder(resp.Body).Decode(&quizData); err != nil {
		return fmt.Errorf("failed to decode quiz data response: %w", err)
	}
	time.Sleep(DefaultWaitTime)
	t.Logger.Info("Quiz started - ID: %d, Number of questions: %d", quizData.ID, quizData.QuestionCount)
	t.Params.Put(questionsKey, quizData.Questions)
	return nil
}
