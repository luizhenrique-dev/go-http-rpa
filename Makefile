default: run-quiz-rpa

.PHONY: install
install:
	@echo "==> Installing application dependencies..."
	@go mod tidy

.PHONY: run-quiz-rpa
run-quiz-rpa:
	@echo "==> Running the Quiz RPA..."
	@go run cmd/quiz/main.go

.PHONY: run-course-rpa
run-course-rpa:
	@echo "==> Running the Course RPA..."
	@go run cmd/watch_course/main.go

.PHONY: run-exam-rpa
run-exam-rpa:
	@echo "==> Running the Answer Exam RPA..."
	@go run cmd/answer_exam/main.go
