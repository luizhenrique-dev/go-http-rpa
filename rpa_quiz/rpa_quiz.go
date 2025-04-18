package rpaquiz

import (
	"os"
	"strings"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	quizPath = "api/quiz/"
	// DefaultWaitTime Default wait time between operations
	DefaultWaitTime = 2 * time.Second
)

// NewRpaQuiz creates a complete job for answering quizzes
func NewRpaQuiz() *engine.Rpa {
	defaultHeaders := make(httprequest.Headers)
	defaultHeaders["Content-Type"] = "application/json"
	defaultHeaders["X-Authorization"] = os.Getenv("WEBSITE_TOKEN")

	quizesIdStr := os.Getenv("WEBSITE_QUIZES_ID")
	var quizIds []string
	if quizesIdStr != "" {
		quizIds = strings.Split(quizesIdStr, ",")
	}

	baseURL := os.Getenv("WEBSITE_URL")
	rpa := engine.NewRpa("rpa_quiz", baseURL, defaultHeaders)
	rpa.SetParams(engine.Parameters{
		"quizIds": quizIds,
		"baseUrl": baseURL,
	})

	rpa.AddTask(NewTaskFetchQuizzes(defaultHeaders, rpa.GetParams()))
	rpa.AddTask(NewTaskProcessQuizes(defaultHeaders, rpa.GetParams()))
	return rpa
}
