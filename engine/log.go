package engine

import (
	"fmt"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, err error, args ...any)
}

type DefaultLogger struct {
	prefix string
}

func (l *DefaultLogger) Info(msg string, args ...any) {
	fmt.Printf("["+l.prefix+"] "+msg+"\n", args...)
}

func (l *DefaultLogger) Error(msg string, err error, args ...any) {
	fmt.Printf("ERROR: "+msg+": %v\n", append(args, err)...)
}
