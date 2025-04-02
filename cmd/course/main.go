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

	quizInput := usecase.CourseInput{
		BaseUrl:   "https://<your-url>/",
		CourseIDs: []int{}, // Add your course ids here. Ex: []int{1, 2, 3}. If not provided it will fetch all available.
		Headers:   headers,
	}
	uc := usecase.NewWatchCourseRpa()
	err := uc.Execute(quizInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RPA executed successfully")
}
