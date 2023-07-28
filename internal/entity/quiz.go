package entity

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
}
