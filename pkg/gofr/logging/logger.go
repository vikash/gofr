package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"cloud.google.com/go/logging"
)

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	level      level
	normalOut  io.Writer
	errorOut   io.Writer
	client     *logging.Client
	isTerminal bool
}

type logEntry struct {
	Level   level       `json:"level"`
	Time    time.Time   `json:"time"`
	Message interface{} `json:"message"`
}

func (l *logger) log(level level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	out := l.normalOut
	if level >= ERROR {
		out = l.errorOut
	}

	entry := logEntry{
		Level: level,
		Time:  time.Now(),
	}

	switch {
	case len(args) == 1 && format == "":
		entry.Message = args[0]
	case len(args) != 1 && format == "":
		entry.Message = args
	case format != "":
		entry.Message = fmt.Sprintf(format+"", args...) // TODO - this is stupid. We should not need empty string.
	}

	if l.isTerminal {
		l.prettyPrint(entry, out)
	} else {
		json.NewEncoder(out).Encode(entry)
	}

}

func (l *logger) Info(args ...interface{}) {
	l.log(INFO, "", args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *logger) Log(args ...interface{}) {
	l.log(INFO, "", args...)
}

func (l *logger) Logf(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log(ERROR, "", args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *logger) prettyPrint(e logEntry, out io.Writer) {
	fmt.Fprintf(out, "\u001B[%dm%s\u001B[0m [%s] %v\n", e.Level.color(), e.Level.String()[0:4], e.Time.Format("15:04:05"), e.Message)
}

func NewLogger(level level) Logger {
	l := &logger{
		normalOut: os.Stdout,
		errorOut:  os.Stderr,
		level:     level,
	}

	client, err := logging.NewClient(context.Background(), "my-project")
	if err != nil {
		l.client = client
	}

	l.isTerminal = checkIfTerminal(l.normalOut)

	return l
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
