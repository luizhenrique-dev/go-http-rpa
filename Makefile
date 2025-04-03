default: run-quiz-rpa

.PHONY: install
install:
	@echo "==> Installing application dependencies..."
	@go mod tidy

.PHONY: run-quiz-rpa
run-quiz-rpa:
	@echo "==> Running Quiz RPA..."
	@go run cmd/quiz/main.go

.PHONY: run-quiz-rpa
run-course-rpa:
	@echo "==> Running Course RPA..."
	@go run cmd/course/main.go
