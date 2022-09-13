package utils

import (
	"fmt"
	"log"
	"monitor-ter/utils"
	"os"
	"path"
	"runtime"
)

var (
	filename = utils.ParamData["LogFile"]
	OK       = fmt.Sprintf("%c[1;0;32m%s%c[0m", 0x1B, "ok", 0x1B)
	WARN     = fmt.Sprintf("%c[1;0;33m%s%c[0m", 0x1B, "warn", 0x1B)
	ERROR    = fmt.Sprintf("%c[1;0;31m%s%c[0m", 0x1B, "error", 0x1B)
	Skip     = fmt.Sprintf("%c[1;0;33m%s%c[0m", 0x1B, "skip", 0x1B)
)

// TerminalLog log
type TerminalLog struct {
}

// GetLogger get logger
func GetLogger(prefix string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	//return log.New(logFile, "["+prefix+"]", log.Ldate|log.Ltime|log.Lshortfile)
	return log.New(logFile, "", log.Ldate|log.Ltime)
}

// InfoFile info
func (ctx TerminalLog) InfoFile(msg string) {
	_, file, lineNo, _ := runtime.Caller(1)
	Log.Printf("[INFO] %v-%v: %v", path.Base(file), lineNo, msg)
}

// SuccessFile success
func (ctx TerminalLog) SuccessFile(msg string) {
	_, file, lineNo, _ := runtime.Caller(1)
	Log.Printf("[SUCCESS] %v-%v: %v", path.Base(file), lineNo, msg)
}

// ErrorFile error
func (ctx TerminalLog) ErrorFile(msg string) {
	_, file, lineNo, _ := runtime.Caller(1)
	Log.Printf("[ERROR] %v-%v: %v", path.Base(file), lineNo, msg)
}

// WarnFile warn
func (ctx TerminalLog) WarnFile(msg string) {
	_, file, lineNo, _ := runtime.Caller(1)

	Log.Printf("[WARN] %v-%v: %v", path.Base(file), lineNo, msg)
}

// SkipFile skip
func (ctx TerminalLog) SkipFile(msg string) {
	_, file, lineNo, _ := runtime.Caller(1)

	Log.Printf("[Skip] %v-%v: %v", path.Base(file), lineNo, msg)
}

var Log = GetLogger("")
var Console = TerminalLog{}
