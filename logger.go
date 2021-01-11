// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

// LogLevel defines the level for log.
// Higher levels log less info.
type LogLevel int

const (
	logPrefix = "rpc"
	//DebugLogLevel defines the level of debug in test environments.
	DebugLogLevel LogLevel = 1
	//TraceLogLevel defines the level of trace in test environments.
	TraceLogLevel LogLevel = 2
	//AllLogLevel defines the lowest level in production environments.
	AllLogLevel LogLevel = 3
	//InfoLogLevel defines the level of info.
	InfoLogLevel LogLevel = 4
	//NoticeLogLevel defines the level of notice.
	NoticeLogLevel LogLevel = 5
	//WarnLogLevel defines the level of warn.
	WarnLogLevel LogLevel = 6
	//ErrorLogLevel defines the level of error.
	ErrorLogLevel LogLevel = 7
	//PanicLogLevel defines the level of panic.
	PanicLogLevel LogLevel = 8
	//FatalLogLevel defines the level of fatal.
	FatalLogLevel LogLevel = 9
	//OffLogLevel defines the level of no log.
	OffLogLevel LogLevel = 10
)
