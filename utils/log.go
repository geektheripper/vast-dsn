package utils

import "log"

type Logger interface {
	Fatalf(format string, v ...interface{})
}

func EnsureLogger(logger ...Logger) Logger {
	if len(logger) > 1 {
		panic("logger must be a single argument")
	}

	if len(logger) == 0 {
		return log.Default()
	}

	return logger[0]
}
