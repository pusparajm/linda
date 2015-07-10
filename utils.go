package main

import (
	log "github.com/Sirupsen/logrus"
)

func stringToLogLevel(str string) log.Level {
	switch str {
	case "panic":
		return log.PanicLevel
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	default:
		return log.ErrorLevel
	}
}
