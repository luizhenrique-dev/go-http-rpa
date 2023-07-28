## GO-HTTP-RPA - Project Execution Guide

This guide provides instructions to execute the project that involves a Go application for automating websites using the **HTTP** protocol.

### Prerequisites

Before you begin, make sure you have the following tools installed on your computer:

- Go (Golang 1.20.6)

### Step 1: Fill in mandatory data

In the **_cmd/quiz/main.go_** file, parameterize the requested data (URL, quiz IDs, authorization token) according to your website.

### Step 2: Run the Go application

In terminal, navigate to the directory where the application project is located. Then, execute the following command to run it:

```
go run cmd/quiz/main.go
```