// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"github.com/hslam/log"
)

// LogLevel defines the level for log.
// Higher levels log less info.
type LogLevel log.Level

const (
	logPrefix = "rpc"
	//AllLogLevel defines the lowest level.
	AllLogLevel = LogLevel(log.AllLevel)
	//TraceLogLevel defines the level of trace in test environments.
	TraceLogLevel = LogLevel(log.TraceLevel)
	//DebugLogLevel defines the level of debug.
	DebugLogLevel = LogLevel(log.DebugLevel)
	//InfoLogLevel defines the level of info.
	InfoLogLevel = LogLevel(log.InfoLevel)
	//NoticeLogLevel defines the level of notice.
	NoticeLogLevel = LogLevel(log.NoticeLevel)
	//WarnLogLevel defines the level of warn.
	WarnLogLevel = LogLevel(log.WarnLevel)
	//ErrorLogLevel defines the level of error.
	ErrorLogLevel = LogLevel(log.ErrorLevel)
	//PanicLogLevel defines the level of panic.
	PanicLogLevel = LogLevel(log.PanicLevel)
	//FatalLogLevel defines the level of fatal.
	FatalLogLevel = LogLevel(log.FatalLevel)
	//OffLogLevel defines the level of no log.
	OffLogLevel = LogLevel(log.OffLevel)
)
