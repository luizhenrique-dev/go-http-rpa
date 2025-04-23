package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	rpaquiz "github.com/luizhenriquees/go-http-rpa/rpa_quiz"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	quizRpa := rpaquiz.NewRpaQuiz()
	if err = quizRpa.Execute(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("RPA executed successfully")
}
