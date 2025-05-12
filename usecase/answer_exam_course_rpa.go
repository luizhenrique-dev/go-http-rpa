package usecase

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/luizhenriquees/go-http-rpa/ai"
	"github.com/luizhenriquees/go-http-rpa/entity"
	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
)

const (
	// API endpoints
	answerPath = "/answer"
	finishPath = "/finish"
	startPath  = "/start"

	// Task status
	statusStarted = "started"
)

// AnswerPayload represents the payload for submitting answers
type AnswerPayload struct {
	Answers []int `json:"answers"`
}

// AnswerExamRpa handles answering exam questions using AI
type AnswerExamRpa struct {
	assistant ai.ExamAssistant
	waitTime  time.Duration
}

// NewAnswerExamRpa creates a new instance of AnswerExamRpa
func NewAnswerExamRpa(assistant ai.ExamAssistant) *AnswerExamRpa {
	return &AnswerExamRpa{
		assistant: assistant,
		waitTime:  2 * time.Second,
	}
}

// Execute processes the exam tasks
func (a *AnswerExamRpa) Execute(input CourseInput) error {
	fmt.Println("AnswerExamRpa initiating...")
	courseList, err := fetchCourseStatus(&input)
	if err != nil {
		return fmt.Errorf("failed to fetch courses status: %w", err)
	}
	filterCoursesBasedOnInput(input, courseList)
	examsList := getExamTasks(courseList)

	for _, examID := range examsList {
		fmt.Printf("Starting exam answering process for exam ID: %d\n", examID)
		// Get exam details
		examTask, err := a.getExamTask(input.BaseUrl, examID, input.Headers)
		if err != nil {
			return fmt.Errorf("failed to get exam task: %w", err)
		}

		if examTask.Type != taskTypeExam {
			return fmt.Errorf("task %d is not an exam type", examID)
		}
		if examTask.Status == statusFinished {
			return fmt.Errorf("task %d is already finished", examID)
		}

		fmt.Printf("Exam found: %s (ID: %d) with %d questions\n",
			examTask.Name, examTask.ID, examTask.QuestionsCount)

		// Start the exam if not already started
		if examTask.Status != statusStarted {
			examTask, err = a.startExam(input.BaseUrl, examID, input.Headers)
			if err != nil {
				return fmt.Errorf("failed to start exam: %w", err)
			}
			fmt.Println("Exam started successfully")
		} else {
			fmt.Println("Exam already started, continuing...")
		}

		// Process and answer questions
		answers, err := a.processQuestions(examTask.Questions)
		if err != nil {
			return fmt.Errorf("error processing questions: %w", err)
		}

		// Submit answers
		if err := a.submitAnswers(input.BaseUrl, examID, answers, input.Headers); err != nil {
			return fmt.Errorf("failed to submit answers: %w", err)
		}

		fmt.Println("Exam completed successfully!")
	}
	fmt.Println("AnswerExamRpa finished!")
	return nil
}

// getExamTask fetches the exam task details
func (a *AnswerExamRpa) getExamTask(baseURL string, examID int, headers map[string]string) (*entity.Task, error) {
	url := baseURL + taskPath + strconv.Itoa(examID)
	fmt.Printf("Fetching exam details from: %s\n", url)

	resp, err := httprequest.DoGet(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	var examTask entity.Task
	if err := json.NewDecoder(resp.Body).Decode(&examTask); err != nil {
		return nil, fmt.Errorf("error decoding exam task: %w", err)
	}

	return &examTask, nil
}

// startExam initiates the exam
func (a *AnswerExamRpa) startExam(baseURL string, examID int, headers map[string]string) (*entity.Task, error) {
	url := baseURL + taskPath + strconv.Itoa(examID) + startPath
	fmt.Printf("Starting exam with request to: %s\n", url)

	resp, err := httprequest.DoPost(url, headers, []byte{})
	if err != nil {
		return nil, fmt.Errorf("error making POST request to start exam: %w", err)
	}
	defer resp.Body.Close()

	var examTask entity.Task
	if err := json.NewDecoder(resp.Body).Decode(&examTask); err != nil {
		return nil, fmt.Errorf("error decoding started exam task: %w", err)
	}

	return &examTask, nil
}

// processQuestions processes each question and gets AI generated answers
func (a *AnswerExamRpa) processQuestions(questions []entity.Question) ([]int, error) {
	answers := make([]int, len(questions))

	for i, question := range questions {
		fmt.Printf("Processing question %d of %d\n", i+1, len(questions))

		// If the question was already answered, use that answer
		if question.Answered != nil {
			fmt.Printf("Question %d was already answered with option %d\n", i+1, *question.Answered)
			answers[i] = *question.Answered
			continue
		}

		// Get answer from AI
		answerIndex, err := a.assistant.GetAnswerIndex(question)
		if err != nil {
			return nil, fmt.Errorf("error getting answer from AI for question %d: %w", i+1, err)
		}

		fmt.Printf("AI selected answer %d for question %d\n", answerIndex, i+1)
		answers[i] = answerIndex

		time.Sleep(a.waitTime)
	}

	return answers, nil
}

// submitAnswers submits the answers to the exam
func (a *AnswerExamRpa) submitAnswers(baseURL string, examID int, answers []int, headers map[string]string) error {
	// First submit answers
	answerURL := baseURL + taskPath + strconv.Itoa(examID) + answerPath
	payload := AnswerPayload{
		Answers: answers,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling answer payload: %w", err)
	}

	fmt.Printf("Submitting answers to: %s\n", answerURL)
	fmt.Printf("Answer payload: %s\n", string(payloadJSON))

	_, err = httprequest.DoPost(answerURL, headers, payloadJSON)
	if err != nil {
		return fmt.Errorf("error submitting answers: %w", err)
	}

	finishURL := baseURL + taskPath + strconv.Itoa(examID) + finishPath
	fmt.Printf("Finishing exam with request to: %s\n", finishURL)

	_, err = httprequest.DoPost(finishURL, headers, payloadJSON)
	if err != nil {
		return fmt.Errorf("error finishing exam: %w", err)
	}

	return nil
}

// getExamTasks extracts all tasks of the type "exam" from the course list
func getExamTasks(courseList *entity.CoursesList) []int {
	if courseList == nil {
		return nil
	}
	var examTasks []int
	for _, course := range courseList.Courses {
		for _, module := range course.Modules {
			for _, task := range module.Tasks {
				if task.Type == taskTypeExam && task.Status != statusFinished {
					examTasks = append(examTasks, task.ID)
					fmt.Printf("Found exam task: %s (ID: %d) in course %d, module %d\n",
						task.Name, task.ID, course.ID, module.ID)
				}
			}
		}
	}
	fmt.Printf("Total exam tasks found: %d\n", len(examTasks))
	return examTasks
}
