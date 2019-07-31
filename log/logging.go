package log

import (
	"log"
	"os"
)

type Level int

const (
	LogPrefix 	= "rpc"
	ALLLevel		Level = -1
	TraceLevel		Level = 1
	DebugLevel		Level = 2
	AllInfoLevel	Level = 3
	InfoLevel		Level = 4
	WarnLevel		Level = 5
	ErrorLevel		Level = 6
	PanicLevel		Level = 7
	FatalLevel		Level = 8
	OffLevel		Level = 9
)

var logLevel Level
var logger *log.Logger

func init() {
	SetLogLevel(Level(6))
	logger = log.New(os.Stdout, "["+LogPrefix+"] ",log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
}

func LogLevel() Level {
	return logLevel
}

func SetLogLevel(level Level) {
	logLevel = level
}


func All(v ...interface{}) {
	logger.Print(v...)
}

func Allf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
func Allln(v ...interface{}) {
	logger.Println(v...)
}


func Trace(v ...interface{}) {
	if logLevel <= TraceLevel {
		logger.Print(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if logLevel <= TraceLevel {
		logger.Printf(format, v...)
	}
}
func Traceln(v ...interface{}) {
	if logLevel <= TraceLevel {
		logger.Println(v...)
	}
}

func Debug(v ...interface{}) {
	if logLevel <= DebugLevel {
		logger.Print(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if logLevel <= DebugLevel {
		logger.Printf(format, v...)
	}
}
func Debugln(v ...interface{}) {
	if logLevel <= DebugLevel {
		logger.Println(v...)
	}
}

func AllInfo(v ...interface{}) {
	if logLevel <= AllInfoLevel {
		logger.Print(v...)
	}
}

func AllInfof(format string, v ...interface{}) {
	if logLevel <= AllInfoLevel {
		logger.Printf(format, v...)
	}
}
func AllInfoln(v ...interface{}) {
	if logLevel <= AllInfoLevel {
		logger.Println(v...)
	}
}
func Info(v ...interface{}) {
	if logLevel <= InfoLevel {
		logger.Print(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		logger.Printf(format, v...)
	}
}
func Infoln(v ...interface{}) {
	if logLevel <= InfoLevel {
		logger.Println(v...)
	}
}

func Warn(v ...interface{}) {
	if logLevel <= WarnLevel {
		logger.Print(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		logger.Printf(format, v...)
	}
}
func Warnln(v ...interface{}) {
	if logLevel <= WarnLevel {
		logger.Println(v...)
	}
}


func Error(v ...interface{}) {
	if logLevel <= ErrorLevel {
		logger.Print(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if logLevel <= ErrorLevel {
		logger.Printf(format, v...)
	}
}
func Errorln(v ...interface{}) {
	if logLevel <= ErrorLevel {
		logger.Println(v...)
	}
}


func Panic(v ...interface{}) {
	if logLevel <= PanicLevel {
		logger.Panic(v...)
	}
}

func Panicf(format string, v ...interface{}) {
	if logLevel <= PanicLevel {
		logger.Panicf(format, v...)
	}
}
func Panicln(v ...interface{}) {
	if logLevel <= PanicLevel {
		logger.Panicln(v...)
	}
}


func Fatal(v ...interface{}) {
	if logLevel <= FatalLevel {
		logger.Fatal(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if logLevel <= FatalLevel {
		logger.Fatalf(format, v...)
	}
}
func Fatalln(v ...interface{}) {
	if logLevel <= FatalLevel {
		logger.Fatalln(v...)
	}
}

func Off(v ...interface{}) {
	if logLevel <= OffLevel {
		logger.Print(v...)
	}
}

func Offf(format string, v ...interface{}) {
	if logLevel <= OffLevel {
		logger.Printf(format, v...)
	}
}
func Offln(v ...interface{}) {
	if logLevel <= OffLevel {
		logger.Println(v...)
	}
}


