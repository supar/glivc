package main

import (
	"github.com/astaxie/beego/logs"
)

// RFC5424 log message levels.
// 0       Emergency: system is unusable
// 1       Alert: action must be taken immediately
// 2       Critical: critical conditions
// 3       Error: error conditions
// 4       Warning: warning conditions
// 5       Notice: normal but significant condition
// 6       Informational: informational messages
// 7       Debug: debug-level messages
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// Wrap external library for the short way to log events.
// This will help to implement other logic
type Log struct {
	// In future its's better to discard chinese code
	*logs.BeeLogger

	// Log level
	Level int
}

// Create log object
func NewLogger(recbuf int64) (logger *Log) {
	logger = &Log{
		BeeLogger: logs.NewLogger(recbuf),
		Level:     LevelDebug,
	}

	// Set default log level
	logger.SetLevel(LevelDebug)

	return
}

// Override to set level value: main class and embeded BeeLog
func (this *Log) SetLevel(lv int) {
	if lv < LevelEmergency || lv > LevelDebug {
		lv = LevelError
	}

	this.Level = lv
	this.BeeLogger.SetLevel(this.Level)
}
