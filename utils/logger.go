package utils

import (
	"log"
	"os"
)

var (
	filename = ParamData["LogFile"]
)

// GetLogger aaa
func GetLogger(prefix string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(logFile, "["+prefix+"]", log.Ldate|log.Ltime|log.Lshortfile)
}

var Log = GetLogger("")
