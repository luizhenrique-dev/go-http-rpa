package main

import (
	"fmt"
	"log"

	"github.com/luizhenriquees/go-http-rpa/internal/usecase"
)

func main() {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["X-Authorization"] = "Bearer <your-token>"

	quizInput := usecase.QuizInput{
		BaseUrl:  "https://<your-url>/",
		QuizesId: []int{}, // Add your quiz ids here. Ex: []int{1, 2, 3}. If not provided it will fetch all available quizzes.
		Headers:  headers,
	}
	uc := usecase.NewAnswerQuizRpa()
	err := uc.Execute(quizInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RPA executed successfully")
}
