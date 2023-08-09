package logger

import (
	"fmt"
	"log"
	"time"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Debug(message string) {
	l.logMessage(logLevelDebug, message)
}

func (l *ConsoleLogger) Info(message string) {
	l.logMessage(logLevelInfo, message)
}

func (l *ConsoleLogger) Error(message string) {
	l.logMessage(logLevelError, message)
}

func (l *ConsoleLogger) logMessage(level string, message string) {
	logMessage := fmt.Sprintf("[%s] %s - %s", level, time.Now().Format(time.RFC3339), message)
	log.Println(logMessage)
}
