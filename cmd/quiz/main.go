package main

import (
	"github.com/luizhenriquees/go-http-rpa/internal/usecase"
)

func main() {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["X-Authorization"] = "Bearer <your-token>"

	quizInput := usecase.QuizInput{
		BaseUrl:  "https://<your-url>/",
		QuizesId: []int{}, // Add your quiz ids here. Ex: []int{1, 2, 3}
		Headers:  headers,
	}
	uc := usecase.NewAnswerQuizRpa()
	uc.Execute(quizInput)
}
