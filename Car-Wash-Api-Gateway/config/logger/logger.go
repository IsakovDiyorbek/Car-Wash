package logger

import (
	"log"
	"os"
)

type Logger struct {
	INFO  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
	DEBUG *log.Logger
	TRACE *log.Logger
}

func NewLogger(basepath, path string) *Logger {
	l := &Logger{}

	fullpath := basepath + "/" + path
	var file, err = os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	l.INFO = log.New(file, "[INFO]  ", log.Lshortfile|log.LstdFlags)
	l.WARN = log.New(file, "[WARN]  ", log.Lshortfile|log.LstdFlags)
	l.ERROR = log.New(file, "[ERROR]  ", log.Lshortfile|log.LstdFlags)
	l.DEBUG = log.New(file, "[DEBUG]  ", log.Lshortfile|log.LstdFlags)
	l.TRACE = log.New(file, "[TRACE]  ", log.Lshortfile|log.LstdFlags)

	return l
}
