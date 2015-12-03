package logger

import (
	"io"
	"io/ioutil"
	"log"
)

//Pre Defined levels log
const (
	LevelNone  = 0
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3
	LevelDebug = 4
	LevelTrace = 5
)

//Loggers by level
var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
	Trace *log.Logger
)

func init() {
	Init(LevelNone, nil, nil, nil, nil, nil)
}

//Init the loggers
func Init(level int, wInfo, wWarn, wError, wDebug, wTrace io.Writer) {
	if level < LevelTrace {
		wTrace = ioutil.Discard
	}
	if level < LevelDebug {
		wDebug = ioutil.Discard
	}
	if level < LevelError {
		wError = ioutil.Discard
	}
	if level < LevelWarn {
		wWarn = ioutil.Discard
	}
	if level < LevelInfo {
		wInfo = ioutil.Discard
	}
	Info = log.New(wInfo, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(wWarn, "WARNING:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(wError, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(wDebug, "DEBUG:", log.Ldate|log.Ltime|log.Lshortfile)
	Trace = log.New(wTrace, "TRACE:", log.Ldate|log.Ltime|log.Lshortfile)
}

//A LevelName pairs of order and name
type LevelName struct {
	Level int
	Name  string
}

//LevelByOrder get Level by order ny [0..3]
func LevelByOrder(ord int) *LevelName {
	switch ord {
	case 0:
		return &LevelName{LevelNone, "None"}
	case 1:
		return &LevelName{LevelInfo, "Info"}
	case 2:
		return &LevelName{LevelWarn, "Warn"}
	default:
		return &LevelName{LevelTrace, "Trace"}
	}
}
