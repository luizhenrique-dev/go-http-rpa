package rpaquiz

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

// TaskFetchQuizzes is a task to fetch all available quizzes
type TaskFetchQuizzes struct {
	*engine.HTTPTask
}

func NewTaskFetchQuizzes(headers httprequest.Headers, params engine.Parameters) *TaskFetchQuizzes {
	return &TaskFetchQuizzes{
		HTTPTask: engine.NewHTTPTask(
			"fetch_quizzes",
			httprequest.GET,
			headers,
			params,
		),
	}
}

func (t *TaskFetchQuizzes) Execute() error {
	quizIds := t.Params.Get("quizIds").([]string)
	if len(quizIds) > 0 {
		return nil
	}
	if err := t.HTTPTask.Execute(); err != nil {
		return err
	}
	var quizList QuizList
	if err := json.NewDecoder(t.Params["response_fetch_quizzes"].(io.Reader)).Decode(&quizList); err != nil {
		return fmt.Errorf("failed to decode quiz list: %w", err)
	}
	for _, quiz := range quizList.Quizzes {
		quizIds = append(quizIds, strconv.Itoa(quiz.ID))
	}
	// TODO check if necessary the put below
	t.Params.Put("quizIds", quizIds)
	log.Printf("Quizzes ID extracted: %v\n", quizIds)
	return nil
}
