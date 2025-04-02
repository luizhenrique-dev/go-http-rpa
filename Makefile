default: run-quiz-rpa

.PHONY: run-quiz-rpa
run-quiz-rpa:
	@echo "==> Running Quiz RPA..."
	@go run cmd/quiz/main.go

.PHONY: run-quiz-rpa
run-course-rpa:
	@echo "==> Running Course RPA..."
	@go run cmd/course/main.go
