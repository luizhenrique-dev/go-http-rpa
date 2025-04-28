package engine

import (
	"fmt"
)

type Logger interface {
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

type DefaultLogger struct {
	prefix string
}

func (l *DefaultLogger) Info(format string, args ...any) {
	fmt.Printf("[%s] INFO: %s\n", l.prefix, fmt.Sprintf(format, args...))
}
func (l *DefaultLogger) Warn(format string, args ...any) {
	fmt.Printf("[%s] WARN: %s\n", l.prefix, fmt.Sprintf(format, args...))
}
func (l *DefaultLogger) Error(format string, args ...any) {
	fmt.Printf("[%s] ERROR: %s\n", l.prefix, fmt.Sprintf(format, args...))
}
