package usecase

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/luizhenriquees/go-http-rpa/engine"
	"github.com/luizhenriquees/go-http-rpa/entity"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	quizPath = "api/quiz/"
)

type QuizInput struct {
	BaseUrl  string
	QuizesId []int
	Headers  map[string]string
}

type AnswerQuizRpa struct {
	logger engine.Logger
}

// NewAnswerQuizRpa this RPA is deprecated
func NewAnswerQuizRpa() *AnswerQuizRpa {
	return &AnswerQuizRpa{
		logger: &engine.DefaultLogger{},
	}
}

func (c *AnswerQuizRpa) SetLogger(logger engine.Logger) {
	c.logger = logger
}

func (c *AnswerQuizRpa) Execute(input QuizInput) error {
	if len(input.QuizesId) == 0 {
		if err := c.fetchAllAvailableQuizzes(&input); err != nil {
			c.logger.Error("Error fetching all quizzes available", err)
			return err
		}
	}
	c.logger.Info("Initiating quiz answering for ids: %v", input.QuizesId)
	for _, quizID := range input.QuizesId {
		quizData, err := c.startQuiz(input.BaseUrl, quizID, input.Headers)
		if err != nil {
			c.logger.Error("Error starting quiz", err, "quizID", quizID)
			return err
		}
		time.Sleep(DefaultWaitTime)
		if err := c.answerQuizQuestions(input.BaseUrl, quizID, quizData, input.Headers); err != nil {
			c.logger.Error("Error answering quiz questions", err, "quizID", quizID)
			return err
		}
	}
	return nil
}

func (c *AnswerQuizRpa) buildQuizURL(baseURL string, quizID int, action string) string {
	url := baseURL + quizPath + strconv.Itoa(quizID)
	if action != "" {
		url += "/" + action
	}
	return url
}

func (c *AnswerQuizRpa) startQuiz(baseURL string, quizID int, headers map[string]string) (*entity.QuizData, error) {
	startURL := c.buildQuizURL(baseURL, quizID, "start")
	c.logger.Info("Starting quiz: POST %s", startURL)
	quizData, err := c.doPostRequest(startURL, headers, nil)
	if err != nil {
		return nil, err
	}
	c.logger.Info("Quiz started - ID: %d, Number of questions: %d", quizData.ID, quizData.QuestionCount)
	return quizData, nil
}

func (c *AnswerQuizRpa) doPostRequest(url string, headers map[string]string, body []byte) (*entity.QuizData, error) {
	resp, err := httprequest.DoPost(url, headers, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP post request failed: %w", err)
	}
	defer resp.Body.Close()
	var responseData entity.QuizData
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &responseData, nil
}

func (c *AnswerQuizRpa) answerQuizQuestions(baseURL string, quizID int, quizData *entity.QuizData, headers map[string]string) error {
	answerURL := c.buildQuizURL(baseURL, quizID, "answer")
	c.logger.Info("Answering questions at: %s", answerURL)
	for index, question := range quizData.Questions {
		payload := c.createAnswerPayload(len(question.Options), index)
		c.logQuestionInfo(question, payload)
		respAnswer, err := c.doPostRequest(answerURL, headers, []byte(payload))
		if err != nil {
			return fmt.Errorf("failed to submit answer for question %d: %w", index, err)
		}
		time.Sleep(DefaultWaitTime)
		c.logger.Info("Question %d result - Correct: %d, Answered: %d",
			index, respAnswer.Questions[index].Correct, respAnswer.Questions[index].Answered)
		c.logger.Info("=====================================================")
	}
	return nil
}

func (c *AnswerQuizRpa) logQuestionInfo(question entity.Question, payload string) {
	c.logger.Info("Question: %s", question.Question)
	c.logger.Info("Number of possible answers: %d", len(question.Options))
	c.logger.Info("Selected answer: %s", payload)
}

func (c *AnswerQuizRpa) createAnswerPayload(optionsCount int, questionIndex int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	answerIndex := rng.Intn(optionsCount)
	return fmt.Sprintf(`{"answer":%d,"question_index":%d}`, answerIndex, questionIndex)
}

func (c *AnswerQuizRpa) fetchAllAvailableQuizzes(input *QuizInput) error {
	c.logger.Info("No specific quiz ID provided. Fetching all available quizzes.")
	quizzesURL := input.BaseUrl + quizPath
	c.logger.Info("GET request to: %s", quizzesURL)
	resp, err := httprequest.DoGet(quizzesURL, input.Headers)
	if err != nil {
		return fmt.Errorf("failed to fetch quizzes: %w", err)
	}
	defer resp.Body.Close()
	var quizList entity.QuizList
	if err := json.NewDecoder(resp.Body).Decode(&quizList); err != nil {
		return fmt.Errorf("failed to decode quiz list: %w", err)
	}
	for _, quiz := range quizList.Quizzes {
		input.QuizesId = append(input.QuizesId, quiz.ID)
	}
	c.logger.Info("Quizzes ID extracted: %v", input.QuizesId)
	return nil
}
