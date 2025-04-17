package rpaquiz

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

// TaskStartQuiz is a task to start a quiz
type TaskStartQuiz struct {
	*engine.HTTPTask
}

func NewTaskStartQuiz(headers httprequest.Headers, params engine.Parameters) *TaskStartQuiz {
	return &TaskStartQuiz{
		HTTPTask: engine.NewHTTPTask(
			"start_quiz",
			httprequest.POST,
			headers,
			params,
		),
	}
}

func (t *TaskStartQuiz) Execute() error {
	baseURL := t.Params["baseUrl"].(string)
	quizID := t.Params[engine.CurrentElement].(string)
	startURL := baseURL + quizPath + quizID + "/start"

	t.Logger.Info("Starting quiz: POST %s", startURL)
	quizData, err := doPostRequest(startURL, t.Headers, nil)
	if err != nil {
		return err
	}
	t.Logger.Info("Quiz started - ID: %d, Number of questions: %d", quizData.ID, quizData.QuestionCount)

	// Store the quiz data for subsequent tasks
	t.Params.Put("quizData", quizData)
	return nil
}

func doPostRequest(url string, headers map[string]string, body []byte) (*QuizData, error) {
	resp, err := httprequest.DoPost(url, headers, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP post request failed: %w", err)
	}
	defer resp.Body.Close()
	var responseData QuizData
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	time.Sleep(DefaultWaitTime)
	return &responseData, nil
}
