package logging

import "log"

type Logger interface {
	Error(msg ...interface{})
	Warn(msg ...interface{})
	Info(msg ...interface{})
	Debug(msg ...interface{})
}

type simpleLoggerImpl struct {
}

func (simpleLoggerImpl) Error(msg ...interface{}) {
	log.Print(msg...)
}

func (simpleLoggerImpl) Warn(msg ...interface{}) {
	log.Print(msg...)
}

func (simpleLoggerImpl) Info(msg ...interface{}) {
	log.Print(msg...)
}

func (simpleLoggerImpl) Debug(msg ...interface{}) {
	log.Print(msg...)
}

func NewLogger(name string) Logger {
	return simpleLoggerImpl{}
}
