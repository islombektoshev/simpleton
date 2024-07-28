package simple

import "log"

// Log You can Overwrite  logger
var Log = func(level LogLevel, msg string) {
	log.Println(msg)
}

type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelDebug
)
