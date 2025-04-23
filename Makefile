default: run-rpa-quiz-dep

.PHONY: install
install:
	@echo "==> Installing application dependencies..."
	@go mod tidy

.PHONY: run-quiz-rpa
run-quiz-rpa:
	@echo "==> Running the Quiz RPA..."
	@go run cmd/quiz/main.go

.PHONY: run-rpa-quiz-dep
run-rpa-quiz-dep:
	@echo "==> Running the Deprecated Quiz RPA..."
	@go run deprecated/cmd/quiz/main.go

.PHONY: run-course-rpa-dep
run-course-rpa-dep:
	@echo "==> Running the Deprecated Course RPA..."
	@go run deprecated/cmd/course/main.go
