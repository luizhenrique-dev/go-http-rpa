package rpa

import (
	"fmt"
	"time"
)

const (
	// DefaultWaitTime Default wait time between operations
	DefaultWaitTime = 2 * time.Second
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, err error, args ...any)
}

type DefaultLogger struct{}

func (l *DefaultLogger) Info(msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}

func (l *DefaultLogger) Error(msg string, err error, args ...any) {
	fmt.Printf(msg+": %v\n", append(args, err)...)
}
