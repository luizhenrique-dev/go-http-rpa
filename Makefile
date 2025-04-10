default: run-rpa-quiz-dep

.PHONY: install
install:
	@echo "==> Installing application dependencies..."
	@go mod tidy

.PHONY: run-quiz-rpa
run-rpa-quiz-dep:
	@echo "==> Running the Deprecated Quiz RPA..."
	@go run deprecated/cmd/quiz/main.go

.PHONY: run-quiz-rpa
run-course-rpa-dep:
	@echo "==> Running the Deprecated Course RPA..."
	@go run deprecated/cmd/course/main.go
