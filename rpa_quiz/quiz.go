package rpaquiz

type QuizInput struct {
	BaseUrl  string
	QuizesId []int
	Headers  map[string]string
}

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Correct  int      `json:"correct"`
	Answered int      `json:"answered"`
}

type QuizData struct {
	ID            int `json:"id"`
	QuestionCount int `json:"questions_count"`
	Questions     []Question
	Status        quizStatus `json:"status"`
}

type QuizList struct {
	Quizzes []QuizData `json:"quiz"`
}
