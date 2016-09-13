package log2

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
)

var (
	logger = &Logger{logs.GetBeeLogger()}
)

// Logger 基于BeegoLog.
type Logger struct {
	*logs.BeeLogger
}

// GetLogger 返回默认的自定义Logger
func GetLogger() *Logger {
	return logger
}

// Output = Info
func (l *Logger) Output(calldepth int, s string) error {
	l.Info(s)
	return nil
}

// Errorf logs a message at error level.
func Errorf(f string, v ...interface{}) {
	logger.Error(f, v)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	logger.Error(fmt.Sprint(v))
}

// Warnf compatibility alias for Warning()
func Warnf(f string, v ...interface{}) {
	logger.Warn(f, v...)
}

// Warn compatibility alias for Warning()
func Warn(v ...interface{}) {
	logger.Warn(fmt.Sprint(v))
}

// Noticef logs a message at notice level.
func Noticef(f string, v ...interface{}) {
	logger.Notice(f, v...)
}

// Notice logs a message at notice level.
func Notice(v ...interface{}) {
	logger.Notice(fmt.Sprint(v))
}

// Infof compatibility alias for Info()
func Infof(f string, v ...interface{}) {
	logger.Info(f, v...)
}

// Info compatibility alias for Info()
func Info(v ...interface{}) {
	logger.Info(fmt.Sprint(v))
}

// Debugf logs a message at debug level.
func Debugf(f string, v ...interface{}) {
	logger.Debug(f, v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	logger.Debug(fmt.Sprint(v))
}

// Fatalln Print Log with Error()  followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	logger.Error(fmt.Sprint(v))
	os.Exit(1)
}

// Fatalf Print Log with Error()  followed by a call to os.Exit(1).
func Fatalf(f string, v ...interface{}) {
	logger.Error(f, v...)
	os.Exit(1)
}

func init() {
	// logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(4)
}
