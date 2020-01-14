package rpc

import (
	"log"
	"os"
)

type Level int

const (
	LogPrefix          = "rpc"
	DebugLevel   Level = 1
	TraceLevel   Level = 2
	AllInfoLevel Level = 3
	InfoLevel    Level = 4
	WarnLevel    Level = 5
	ErrorLevel   Level = 6
	PanicLevel   Level = 7
	FatalLevel   Level = 8
	OffLevel     Level = 9
	ALLLevel     Level = 10
	NoLevel      Level = 99
)

var logLevel Level
var Logger *log.Logger

func init() {
	SetLogLevel(InfoLevel)
	InitLog()
}
func InitLog() {
	Logger = log.New(os.Stdout, "["+LogPrefix+"] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
}
func LogLevel() Level {
	return logLevel
}

func SetLogLevel(level Level) {
	logLevel = level
}

func All(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= ALLLevel {
		Logger.Print(v...)
	}
}

func Allf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= ALLLevel {
		Logger.Printf(format, v...)
	}
}
func Allln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= ALLLevel {
		Logger.Println(v...)
	}
}

func Debug(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= DebugLevel {
		Logger.Print(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= DebugLevel {
		Logger.Printf(format, v...)
	}
}
func Debugln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= DebugLevel {
		Logger.Println(v...)
	}
}

func Trace(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= TraceLevel {
		Logger.Print(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= TraceLevel {
		Logger.Printf(format, v...)
	}
}
func Traceln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= TraceLevel {
		Logger.Println(v...)
	}
}
func AllInfo(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= AllInfoLevel {
		Logger.Print(v...)
	}
}

func AllInfof(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= AllInfoLevel {
		Logger.Printf(format, v...)
	}
}
func AllInfoln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= AllInfoLevel {
		Logger.Println(v...)
	}
}
func Info(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= InfoLevel {
		Logger.Print(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= InfoLevel {
		Logger.Printf(format, v...)
	}
}
func Infoln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= InfoLevel {
		Logger.Println(v...)
	}
}

func Warn(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= WarnLevel {
		Logger.Print(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= InfoLevel {
		Logger.Printf(format, v...)
	}
}
func Warnln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= WarnLevel {
		Logger.Println(v...)
	}
}

func Error(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= ErrorLevel {
		Logger.Print(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= ErrorLevel {
		Logger.Printf(format, v...)
	}
}
func Errorln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= ErrorLevel {
		Logger.Println(v...)
	}
}

func Panic(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= PanicLevel {
		Logger.Panic(v...)
	}
}

func Panicf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= PanicLevel {
		Logger.Panicf(format, v...)
	}
}
func Panicln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= PanicLevel {
		Logger.Panicln(v...)
	}
}

func Fatal(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= FatalLevel {
		Logger.Fatal(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= FatalLevel {
		Logger.Fatalf(format, v...)
	}
}
func Fatalln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= FatalLevel {
		Logger.Fatalln(v...)
	}
}

func Off(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= OffLevel {
		Logger.Print(v...)
	}
}

func Offf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= OffLevel {
		Logger.Printf(format, v...)
	}
}
func Offln(v ...interface{}) {
	if Logger == nil {
		return
	}
	if logLevel <= OffLevel {
		Logger.Println(v...)
	}
}
