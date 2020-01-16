package rpc

import (
	"github.com/hslam/log"
)

// LogLevel defines the level for log.
// Higher levels log less info.
type LogLevel int

const (
	logPrefix = "rpc"
	//DebugLevel defines the level of debug in test environments.
	DebugLevel LogLevel = 1
	//TraceLevel defines the level of trace in test environments.
	TraceLevel LogLevel = 2
	//AllLevel defines the lowest level in production environments.
	AllLevel LogLevel = 3
	//InfoLevel defines the level of info.
	InfoLevel LogLevel = 4
	//NoticeLevel defines the level of notice.
	NoticeLevel LogLevel = 5
	//WarnLevel defines the level of warn.
	WarnLevel LogLevel = 6
	//ErrorLevel defines the level of error.
	ErrorLevel LogLevel = 7
	//PanicLevel defines the level of panic.
	PanicLevel LogLevel = 8
	//FatalLevel defines the level of fatal.
	FatalLevel LogLevel = 9
	//OffLevel defines the level of no log.
	OffLevel LogLevel = 10
)

var logger = log.New()

func init() {
	logger.SetPrefix(logPrefix)
	SetLogLevel(InfoLevel)
}

func SetLogLevel(level LogLevel) {
	logger.SetLevel(log.Level(level))
}

func GetLogLevel() LogLevel {
	return LogLevel(logger.GetLevel())
}
