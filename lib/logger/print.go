package logger

import (
	"log"
	"os"
)

/*
#######################
#		log日志工具	   #
#######################
*/

var stdoutLogger = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var stderrLogger = log.New(os.Stderr, "[ERROR]", log.LstdFlags)

func Errorf(format string, v ...interface{}) {
	stderrLogger.Printf(format, v)
}
func Error(v ...interface{}) {
	stderrLogger.Println(v)
}
func Err(err error) {
	stderrLogger.Println(err.Error())
}
func Infof(format string, v ...interface{}) {
	stdoutLogger.Printf(format, v)
}
func Info(v ...interface{}) {
	stdoutLogger.Println(v)
}
