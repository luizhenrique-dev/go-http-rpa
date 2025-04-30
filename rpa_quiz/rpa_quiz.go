package rpaquiz

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	maxPerExec = "maxPerExec"
	quizPath   = "api/quiz/"
	// DefaultWaitTime Default wait time between operations
	DefaultWaitTime = 2 * time.Second
)

// NewRpaQuiz creates a complete job for answering quizzes
func NewRpaQuiz() *engine.Rpa {
	defaultHeaders := make(httprequest.Headers)
	defaultHeaders["Content-Type"] = "application/json"
	defaultHeaders["X-Authorization"] = os.Getenv("WEBSITE_TOKEN")

	maxPerExecution := getMaxPerExecution()
	todoQuizIds := buildQuizIdList(maxPerExecution)
	baseURL := os.Getenv("WEBSITE_URL")
	rpa := engine.NewRpa("rpa_quiz", baseURL, defaultHeaders)
	rpa.SetParams(engine.Parameters{
		quizIdsKey:          todoQuizIds,
		engine.ParamBaseURL: baseURL,
		maxPerExec:          maxPerExecution,
	})

	rpa.AddTask(NewTaskFetchQuizzes(defaultHeaders, rpa.GetParams()))
	rpa.AddTask(NewTaskProcessQuizes(defaultHeaders, rpa.GetParams()))
	return rpa
}

func buildQuizIdList(maxPerExecution int) []string {
	quizesIdStr := os.Getenv("WEBSITE_QUIZES_ID")
	var todoQuizIds []string
	if quizesIdStr != "" {
		requestedQuizIds := strings.Split(quizesIdStr, ",")
		log.Printf("Requested Quiz IDs list: %v\n", requestedQuizIds)
		for i := 0; i < maxPerExecution; i++ {
			if len(requestedQuizIds) < i {
				break
			}
			todoQuizIds = append(todoQuizIds, requestedQuizIds[i])
		}
	}
	log.Printf("Quiz IDs to do list: %v\n", todoQuizIds)
	return todoQuizIds
}

func getMaxPerExecution() int {
	maxPerExecution := os.Getenv("MAX_PER_EXECUTION")
	if maxPerExecution == "" {
		return 0
	}
	maxPerExecutionInt, err := strconv.Atoi(maxPerExecution)
	if err != nil {
		log.Fatal(err)
	}
	return maxPerExecutionInt
}
