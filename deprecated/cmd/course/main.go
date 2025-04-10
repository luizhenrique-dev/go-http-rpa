package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/luizhenriquees/go-http-rpa/deprecated/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["X-Authorization"] = os.Getenv("WEBSITE_TOKEN")

	quizInput := usecase.CourseInput{
		BaseUrl:   os.Getenv("WEBSITE_URL"),
		CourseIDs: []int{}, // Add your course ids here. Ex: []int{1, 2, 3}. If not provided it will fetch all available.
		Headers:   headers,
	}
	uc := usecase.NewWatchCourseRpa()
	err = uc.Execute(quizInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RPA executed successfully")
}
