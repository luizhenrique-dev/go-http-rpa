package usecase

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/luizhenriquees/go-http-rpa/internal/entity"
	http_request "github.com/luizhenriquees/go-http-rpa/pkg/http_request"
)

type QuizInput struct {
	BaseUrl  string
	QuizesId []int
	Headers  map[string]string
}

type AnswerQuizRpa struct {
}

func NewAnswerQuizRpa() *AnswerQuizRpa {
	return &AnswerQuizRpa{}
}

func (c *AnswerQuizRpa) Execute(input QuizInput) error {
	for _, actualQuizId := range input.QuizesId {
		urlStartQuiz := input.BaseUrl + "api/quiz/" + strconv.Itoa(actualQuizId) + "/start"
		fmt.Println("POST start:" + urlStartQuiz)
		resp, err := doPostParsingData(urlStartQuiz, input.Headers, nil)
		if err != nil {
			fmt.Println("Error in POST request 'start':", err)
			return err
		}
		fmt.Println("Response POST:")
		printResponse(resp)
		time.Sleep(2 * time.Second)
		answerEachQuestion(input.BaseUrl, actualQuizId, resp, input.Headers)
	}
	return nil
}

func doPostParsingData(url string, headers map[string]string, body []byte) (*entity.QuizData, error) {
	resp, err := http_request.DoPost(url, headers, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseData entity.QuizData
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func printResponse(resp *entity.QuizData) {
	fmt.Println("ID:", resp.ID)
	fmt.Println("Number of questions:", resp.QuestionCount)
}

func answerEachQuestion(baseUrl string, actualQuizIndex int, resp *entity.QuizData, headers map[string]string) {
	urlQuizAnswer := baseUrl + "api/quiz/" + strconv.Itoa(actualQuizIndex) + "/answer"
	fmt.Println("POST answer: " + urlQuizAnswer)
	for index, question := range resp.Questions {
		payloadAnswer := getPayloadAnswer(len(question.Options), index)
		printQuestionData(question, payloadAnswer)
		respAnswer, err := doPostParsingData(urlQuizAnswer, headers, []byte(payloadAnswer))
		time.Sleep(2 * time.Second)
		if err != nil {
			fmt.Println("Error in POST request 'answer':", err)
			return
		}
		fmt.Printf("RESULT Correct: %d, Answered: %d\n", respAnswer.Questions[index].Correct, respAnswer.Questions[index].Answered)
		fmt.Println("=====================================================")
	}
}

func printQuestionData(question entity.Question, payloadAnswer string) {
	fmt.Printf("Question: %s\n", question.Question)
	fmt.Printf("Number of possible answers: %d\n", len(question.Options))
	fmt.Println("Answered:", payloadAnswer)
}

func getPayloadAnswer(optionsLength int, questionIndex int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randAnswerIndex := random.Intn(optionsLength)
	return `{"answer":` + strconv.Itoa(randAnswerIndex) + `,"question_index":` + strconv.Itoa(questionIndex) + `}`
}
