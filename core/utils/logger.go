package utils

import (
	"log"
)

type ZuubyLogger log.Logger

func (l *ZuubyLogger) Error(message string) {
	l.Println("[error] " + message)
}

func (l *ZuubyLogger) Warn(message string) {
	l.Println("[warn] " + message)
}

func (l *ZuubyLogger) Debug(message string) {
	l.Println("[debug] " + message)
}

func (l *ZuubyLogger) Info(message string) {
	l.Println("[info] " + message)
}

func NewLogger(name string) *ZuubyLogger {
	// TODO: add log rotation
	return log.New(os.Stdout, name, log.Ltime|log.Lshortfile|log.LUTC)
}
