# GO-HTTP-RPA

## Overview
GO-HTTP-RPA is a Robotic Process Automation (RPA) tool built in Go that automates interactions with web-based systems using the HTTP protocol. The project provides automation for:
- Watching online courses and completing tasks
- Answering quizzes automatically

This approach offers a lightweight alternative to browser-based automation, using direct HTTP requests to simulate user interactions.
## Features
- **Course Automation**: Automatically progress through course modules and complete tasks
- **Quiz Automation**: Submit answers to quizzes with random selection
- **Flexible Configuration**: Filter specific courses or quizzes by ID
- **Customizable Timing**: Adjustable wait times between operations

## Prerequisites
- Go 1.22.x or higher
- Basic knowledge of the target system's API endpoints

## Getting Started
### 1. Clone the Repository
``` sh
git clone https://github.com/luizhenrique-dev/go-http-rpa.git
cd go-http-rpa
```

### 2. Configure Your Target System

#### Configure Environment Variables
Create a `.env` file in the project root with the following variables:
``` 
WEBSITE_URL=https://<your-url>/
WEBSITE_TOKEN=Bearer <your-token>
```

#### For Quiz Automation
Edit `cmd/quiz/main.go` and update the following:
``` go
quizInput := usecase.QuizInput{
    QuizesId: []int{}, // Add your quiz ids here. Ex: []int{1, 2, 3}
}
```
#### For Course Automation
Edit `cmd/course/main.go` and update the following:
``` go
courseInput := usecase.CourseInput{
    CourseIDs: []int{}, // Add your course ids here. Ex: []int{1, 2, 3}
}
```
### 3. Run the Application
You can run the application using either direct Go commands or Makefile targets:
#### Using Makefile:
For quiz automation:
``` sh
make run-quiz-rpa
```
For course watching:
``` sh
make run-course-rpa
```
## Project Structure
``` 
go-http-rpa/
├── cmd/
│   ├── quiz/           # Quiz automation entry point
│   └── course/         # Course automation entry point
├── internal/
│   ├── entity/         # Data models and structures
│   └── usecase/        # Business logic implementation
├── pkg/
│   └── http_request/   # HTTP request utilities
└── Makefile            # Build and run commands
```
## How It Works
1. **Authentication**: Uses provided tokens for API authentication
2. **Data Fetching**: Retrieves available courses/quizzes from the system
3. **Task Processing**: Systematically works through modules and tasks
4. **Automated Responses**: Submits randomly selected answers for quizzes and tests

## Configuration Options
### Quiz Automation
- `BaseUrl`: Target system base URL
- `QuizesId`: Array of specific quiz IDs to process (empty array fetches all)
- `Headers`: HTTP headers for authentication and content type

### Course Automation
- `BaseUrl`: Target system base URL
- `CourseIDs`: Array of specific course IDs to process (empty array processes all)
- `Headers`: HTTP headers for authentication and content type
- `WithWaitTime`: Optional parameter to customize wait times between operations
