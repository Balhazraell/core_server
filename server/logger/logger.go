package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

var (
	INFO    = "INFO: "
	WARNING = "WARNING: "
	ERROR   = "ERROR: "
)

var (
	file          io.Writer
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func InitLogger() bool {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open or create log file: ", err)
		return false
	}

	infoLogger = log.New(
		file,
		INFO,
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	warningLogger = log.New(
		file,
		WARNING,
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	errorLogger = log.New(
		file,
		ERROR,
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	return true
}

func getCallStackLine(depth ...int) string {
	var currentDepth = 2

	if len(depth) > 0 {
		currentDepth = depth[0]
	}

	_, callFile, line, ok := runtime.Caller(currentDepth)

	if !ok {
		callFile = "???"
		line = 0
	}

	return fmt.Sprintf("%s:%v:", callFile, line)
}

//------------------- Info -------------------//
func InfoPrint(a ...interface{}) {
	fmt.Println(INFO, getCallStackLine(), fmt.Sprint(a...))
	infoLogger.Println(getCallStackLine(), fmt.Sprint(a...))
}

func InfoPrintf(str string, a ...interface{}) {
	fmt.Println(INFO, getCallStackLine(), fmt.Sprintf(str, a...))
	infoLogger.Println(getCallStackLine(), fmt.Sprintf(str, a...))
}

//------------------- Warning -------------------//
func WarningPrint(a ...interface{}) {
	fmt.Println(WARNING, getCallStackLine(), fmt.Sprint(a...))
	warningLogger.Println(getCallStackLine(), fmt.Sprint(a...))
}

func WarningPrintf(str string, a ...interface{}) {
	fmt.Println(WARNING, getCallStackLine(), fmt.Sprintf(str, a...))
	warningLogger.Println(getCallStackLine(), fmt.Sprintf(str, a...))
}

//------------------- Error -------------------//
func ErrorPrint(a ...interface{}) {
	fmt.Println(ERROR, getCallStackLine(), fmt.Sprint(a...))
	errorLogger.Println(getCallStackLine(), fmt.Sprint(a...))
}

func ErrorPrintf(str string, a ...interface{}) {
	fmt.Println(ERROR, getCallStackLine(), fmt.Sprintf(str, a...))
	errorLogger.Println(getCallStackLine(), fmt.Sprintf(str, a...))
}
