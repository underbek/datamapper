package logger

import (
	"fmt"
	"log"
	"os"
)

const depth = 2

type Logger interface {
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
}

type logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func New() Logger {
	return &logger{
		info:  log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime),
		warn:  log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
	}
}

func (l *logger) Info(v ...any) {
	_ = l.info.Output(depth, fmt.Sprint(v...))
}

func (l *logger) Infof(format string, v ...any) {
	_ = l.info.Output(depth, fmt.Sprintf(format, v...))
}

func (l *logger) Warn(v ...any) {
	_ = l.warn.Output(depth, fmt.Sprint(v...))
}

func (l *logger) Error(v ...any) {
	_ = l.error.Output(depth, fmt.Sprint(v...))
}

func (l *logger) Errorf(format string, v ...any) {
	_ = l.error.Output(depth, fmt.Sprintf(format, v...))
}

func (l *logger) Fatal(v ...any) {
	_ = l.error.Output(depth, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...any) {
	_ = l.error.Output(depth, fmt.Sprintf(format, v...))
	os.Exit(1)
}
