package logger

import (
	"log"
	"os"
)

type LoggerInterface interface {
	Info(msg string)
	Error(msg string)
}

type logger struct {
	info  *log.Logger
	error *log.Logger
}

func NewLogger() LoggerInterface {
	return &logger{
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *logger) Info(msg string) {
	l.info.Println(msg)
}

func (l *logger) Error(msg string) {
	l.error.Println(msg)
}
