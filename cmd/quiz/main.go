package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/luizhenriquees/go-http-rpa/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["X-Authorization"] = os.Getenv("WEBSITE_TOKEN")

	quizInput := usecase.QuizInput{
		BaseUrl:  os.Getenv("WEBSITE_URL"),
		QuizesId: []int{}, // Add your quiz ids here. Ex: []int{1, 2, 3}. If not provided it will fetch all available quizzes.
		Headers:  headers,
	}
	uc := usecase.NewAnswerQuizRpa()
	err = uc.Execute(quizInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RPA executed successfully")
}
