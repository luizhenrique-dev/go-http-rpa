package rpaquiz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/luizhenriquees/go-http-rpa/engine"
	"github.com/luizhenriquees/go-http-rpa/entity"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	fetchQuizPath string            = "api/quiz"
	pending       entity.QuizStatus = "pending"
)

// TaskFetchQuizzes is a task to fetch all available quizzes
type TaskFetchQuizzes struct {
	*engine.HTTPTask
}

func NewTaskFetchQuizzes(headers httprequest.Headers, params engine.Parameters) *TaskFetchQuizzes {
	task := &TaskFetchQuizzes{}
	httpTask := engine.NewHTTPTask(
		"fetch_quizzes",
		httprequest.GET,
		headers,
		params,
		engine.WithPreRequestFunc(task.preRequest),
		engine.WithPostExtractFunc(task.postExtract),
	)
	task.HTTPTask = httpTask
	return task
}

func (t *TaskFetchQuizzes) preRequest() error {
	baseURL := ""
	if val, ok := t.Params.Get(engine.ParamBaseURL).(string); ok {
		baseURL = val
	}
	fetchQuizURL := baseURL + fetchQuizPath
	t.Logger.Info("pre-request built URL: %s", fetchQuizURL)
	t.URL = fetchQuizURL
	return nil
}

func (t *TaskFetchQuizzes) postExtract(resp *http.Response, _ *engine.HTTPTask) error {
	t.Logger.Info("PostExtract used.")
	defer resp.Body.Close()
	maxPerExecution := t.Params.Get(maxPerExec).(int)
	pendingQuizIds := t.Params.Get(quizIdsKey).([]string)
	if len(pendingQuizIds) > 0 {
		return nil
	}
	var quizList entity.QuizList
	if err := json.NewDecoder(resp.Body).Decode(&quizList); err != nil {
		return fmt.Errorf("failed to decode quiz list: %w", err)
	}
	for _, quiz := range quizList.Quizzes {
		if quiz.Status != pending {
			continue
		}
		if len(pendingQuizIds) < maxPerExecution {
			pendingQuizIds = append(pendingQuizIds, strconv.Itoa(quiz.ID))
		}
	}
	t.Params.Put(quizIdsKey, pendingQuizIds)
	log.Printf("Quiz IDs to do list: %v\n", pendingQuizIds)
	return nil
}
