package usecase

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	httprequest "github.com/luizhenriquees/go-http-rpa/http_request"
	"github.com/luizhenriquees/go-http-rpa/rpa_course"
)

const (
	// API endpoints
	taskPath   = "api/task/"
	statusPath = "api/status"

	// Task types
	taskTypeExam = "exam"
	taskTypeTest = "test"

	// Task status
	statusFinished  = "finished"
	DefaultWaitTime = 2 * time.Second
)

type CourseInput struct {
	BaseUrl   string
	CourseIDs []int
	Headers   map[string]string
}

type WatchCourseRpa struct {
	waitTime time.Duration
}

type WatchCourseOption func(*WatchCourseRpa)

// NewWatchCourseRpa this RPA is deprecated
func NewWatchCourseRpa(opts ...WatchCourseOption) *WatchCourseRpa {
	rpa := &WatchCourseRpa{
		waitTime: DefaultWaitTime,
	}
	for _, opt := range opts {
		opt(rpa)
	}
	return rpa
}

func (w *WatchCourseRpa) Execute(input CourseInput) error {
	courseList, err := fetchCourseStatus(&input)
	if err != nil {
		return fmt.Errorf("failed to fetch courses status: %w", err)
	}
	filterCoursesBasedOnInput(input, courseList)
	for _, course := range courseList.Courses {
		if err := w.processCourse(input, course); err != nil {
			return fmt.Errorf("error processing course %d: %w", course.ID, err)
		}
	}
	return nil
}

func (w *WatchCourseRpa) processCourse(input CourseInput, course rpacourse.Course) error {
	fmt.Printf("Watching Course ID: %d\n", course.ID)
	for _, module := range course.Modules {
		if err := w.processModule(input, course.ID, module); err != nil {
			return fmt.Errorf("error processing module %d: %w", module.ID, err)
		}
	}
	return nil
}

func (w *WatchCourseRpa) processModule(input CourseInput, courseID int, module rpacourse.Module) error {
	fmt.Printf("Watching Module ID: %d\n", module.ID)
	for _, task := range module.Tasks {
		if task.Type == taskTypeExam {
			fmt.Printf("[Course %d] | [Module %d] - Task %d is an exam! Stopping...\n",
				courseID, module.ID, task.ID)
			break
		}
		if err := w.processTask(input, courseID, module.ID, task); err != nil {
			return fmt.Errorf("error processing task %d: %w", task.ID, err)
		}
		time.Sleep(w.waitTime)
	}
	return nil
}

func (w *WatchCourseRpa) processTask(input CourseInput, courseID, moduleID int, task rpacourse.Task) error {
	startedTask, err := w.startTask(input, courseID, moduleID, task)
	if err != nil {
		return err
	}
	time.Sleep(w.waitTime)
	var questionAnsweredBody []byte
	if w.isTaskATest(startedTask) {
		fmt.Printf("Task %d is a single test! Building answer...\n", task.ID)
		answerJSON := w.buildCourseTestAnswer(len(startedTask.Questions[0].Options))
		questionAnsweredBody = []byte(answerJSON)
	}
	return w.finishTask(input, courseID, moduleID, task.ID, questionAnsweredBody)
}

func (w *WatchCourseRpa) isTaskATest(startedTask *rpacourse.Task) bool {
	return startedTask.Type == taskTypeTest && startedTask.QuestionsCount == 1 && startedTask.Status != statusFinished
}

func (w *WatchCourseRpa) startTask(input CourseInput, courseID, moduleID int, task rpacourse.Task) (*rpacourse.Task, error) {
	urlStartTask := input.BaseUrl + taskPath + strconv.Itoa(task.ID) + "/start"
	respStartTask, err := httprequest.DoPost(urlStartTask, input.Headers, []byte(""))
	if err != nil {
		return nil, fmt.Errorf("error starting task %d: %w", task.ID, err)
	}
	defer respStartTask.Body.Close()
	fmt.Printf("[Course %d] | [Module %d] - Task %d started!\n", courseID, moduleID, task.ID)
	var startedTask rpacourse.Task
	if err := json.NewDecoder(respStartTask.Body).Decode(&startedTask); err != nil {
		return nil, fmt.Errorf("error parsing started task: %w", err)
	}
	return &startedTask, nil
}

func (w *WatchCourseRpa) finishTask(input CourseInput, courseID, moduleID, taskID int, answerBody []byte) error {
	urlFinishTask := input.BaseUrl + taskPath + strconv.Itoa(taskID) + "/finish"
	_, err := httprequest.DoPost(urlFinishTask, input.Headers, answerBody)
	if err != nil {
		return fmt.Errorf("error finishing task %d: %w", taskID, err)
	}
	fmt.Printf("[Course %d] | [Module %d] - Task %d finished!\n", courseID, moduleID, taskID)
	return nil
}

func (w *WatchCourseRpa) buildCourseTestAnswer(optionsLength int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randAnswerIndex := random.Intn(optionsLength)
	fmt.Printf("Answer index chosen: %d\n", randAnswerIndex)
	return fmt.Sprintf(`{"answers":[%d]}`, randAnswerIndex)
}

func fetchCourseStatus(input *CourseInput) (*rpacourse.CoursesList, error) {
	urlGetCourses := input.BaseUrl + statusPath
	fmt.Println("GET to:", urlGetCourses)
	resp, err := httprequest.DoGet(urlGetCourses, input.Headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching courses list: %w", err)
	}
	defer resp.Body.Close()
	var responseData rpacourse.CoursesList
	fmt.Println("Courses fetched, extracting data...")
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("error decoding courses list: %w", err)
	}
	fmt.Printf("%d courses extracted...\n", len(responseData.Courses))
	return &responseData, nil
}

func filterCoursesBasedOnInput(input CourseInput, courseList *rpacourse.CoursesList) {
	if len(input.CourseIDs) == 0 {
		fmt.Println("No specific course ID provided. All available courses will be watched.")
		return
	}
	fmt.Printf("Filtering course IDs: %v\n", input.CourseIDs)
	idMap := make(map[int]bool)
	for _, id := range input.CourseIDs {
		idMap[id] = true
	}
	var filteredCourses []rpacourse.Course
	for _, course := range courseList.Courses {
		if idMap[course.ID] && course.Status != statusFinished {
			filteredCourses = append(filteredCourses, course)
		}
	}
	courseList.Courses = filteredCourses
	fmt.Printf("%d courses remaining after filter...\n", len(courseList.Courses))
}
