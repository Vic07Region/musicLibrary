package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
)

// Цвета для уровней логирования
var levelColors = map[int]string{
	InfoLevel:  "\033[1;36m", // Cyan
	DebugLevel: "\033[1;33m", // Orange
	WarnLevel:  "\033[1;31m", // Orange
}

type Logger struct {
}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logMessage(DebugLevel, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logMessage(InfoLevel, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logMessage(WarnLevel, msg, args...)
}

func (l *Logger) Fatal(v ...any) {
	log.Fatal(v)
}

func (l *Logger) logMessage(level int, msg string, args ...interface{}) {
	data := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			data[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			data[fmt.Sprintf("%v", args[i])] = ""
		}
	}

	jsonData, _ := json.MarshalIndent(data, " ", "  ")
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	formattedMessage := fmt.Sprintf(
		"%s %s %s %s",
		timestamp,
		levelColors[level],
		l.levelToString(level),
		msg,
	)
	if len(args) > 0 {
		fmt.Printf(formattedMessage+"\u001B[0;35m %s\u001B[0m\n", jsonData)
	} else {
		fmt.Println(formattedMessage + "\u001B[0m")
	}
}

func (l *Logger) levelToString(level int) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	default:
		return "UNKNOWN"
	}
}
