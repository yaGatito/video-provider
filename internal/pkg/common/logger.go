package common

import (
	"io"
	"log"
	"os"
	"time"
)

// TODO: reconsider using sb
var DefaultOutput = os.Stdout

type Logger struct {
	svcID string
	log   *log.Logger
}

func NewLogger(output io.Writer, svcID string) *Logger {
	return &Logger{
		svcID: svcID,
		log:   log.New(output, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC),
	}
}

func (l *Logger) Info(message string) {
	l.log.Printf("%s [INFO] %s: %s\n", time.Now().String(), l.svcID, message)
}

func (l *Logger) Debug(message string) {
	l.log.Printf("%s [DEBUG] %s: %s\n", time.Now().String(), l.svcID, message)
}

func (l *Logger) Error(message string, err error) {
	l.log.Printf("%s [ERROR] %s: %s\n", time.Now().String(), l.svcID, message)
	if err != nil {
		l.Debug(err.Error())
	}
}
