package arangogo

import "log"

type Logger interface {
	Log(msg string)
}

type stdLoggerBasedLogger struct {
	logger *log.Logger
}

func (l stdLoggerBasedLogger) Log(msg string) {
	l.logger.Println(msg)
}

func NewLoggerWithStdLogger(logger *log.Logger) Logger {
	return &stdLoggerBasedLogger{logger: logger}
}
