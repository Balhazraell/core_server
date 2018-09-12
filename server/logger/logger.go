package logger

import (
	"io"
	"log"
	"os"
)

var (
	file    io.Writer
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Идея такая, инициализацию логгера инициирует core часть при запуске.
// Параметры запуска диктует файл с настройками.
func InitLogger() bool {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open or create log file: ", err)
		return false
	}

	Info = log.New(
		file,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Warning = log.New(
		file,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Error = log.New(
		file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	return true
}
