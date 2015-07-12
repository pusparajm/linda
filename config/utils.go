package config

import (
	log "github.com/Sirupsen/logrus"
)

// Returns built-in Help command
func NewHelpCommand() Command {
	return Command{
		Type:        "Snitch",
		Name:        "help",
		Description: "Displays a list of commands. Usage: `!help`",
		Expression:  "(?i)!help",
	}
}

// Returns execution mode by input
func GetExecutionMode(input string) string {
	switch input {
	case ExecutionModeFirst:
		return ExecutionModeFirst
	case ExecutionModeAll:
		return ExecutionModeAll
	default:
		return ExecutionModeFirst
	}
}

// Returns logrus.LogLevel by provided string
func StringToLogLevel(str string) log.Level {
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
