package mylog

import (
	"io"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// MyLogger interface
type MyLogger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	GetEntryWithFields(m map[string]interface{}) *Entry
	GetField(key string) interface{}
}

// Logger log interface
type Logger struct {
	*log.Logger

	welog *log.Logger // logger for warning and error
	mu    sync.Mutex
	preM  log.Fields
}

const (
	// Lfile log with fileline
	Lfile = 1 << iota
	// Ljson json format
	Ljson
	// Lrelative use relative path
	Lrelative
	// Lwelog copy warning and error log to another output
	Lwelog
)

// LogLevel map level string to log.Level
var LogLevel = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
	"WARN":  log.WarnLevel,
	"ERROR": log.ErrorLevel,
	"FATAL": log.FatalLevel,
}

var stdLog = New("DEBUG", os.Stderr, Lfile)

func getLogLevel(level string) log.Level {
	l, ok := LogLevel[strings.ToUpper(level)]
	if !ok {
		return log.DebugLevel
	}

	return l
}

// New return a new Logger, Lrelative, Lwelog is invalid in this function
func New(level string, out io.Writer, flag int) *Logger {
	return NewWithRelativePath(level, out, flag, "", nil)
}

// NewWithRelativePath return a new Logger, relativePath: ignore prefix from path
func NewWithRelativePath(level string, out io.Writer, flag int, relativePath string, errout io.Writer) *Logger {
	myLogger := &Logger{Logger: newLogger(level, out, flag, relativePath)}
	if flag&Lwelog != 0 && errout != nil {
		myLogger.welog = newLogger(level, errout, flag, relativePath)
	}

	return myLogger
}

func newLogger(level string, out io.Writer, flag int, relativePath string) *log.Logger {
	Log := log.New()
	Log.Level = getLogLevel(level)
	Log.Out = out
	if flag&Ljson != 0 {
		Log.Formatter = &log.JSONFormatter{}
	}

	if flag&Lfile != 0 {
		flHook := FileLineHook{}
		if flag&Lrelative != 0 {
			flHook.IsRelative = true
			flHook.IgnorePrefix = relativePath
		}

		Log.Hooks.Add(flHook)
	}

	return Log
}

// Predefine log with fixed extra record
func (l *Logger) Predefine(m map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.preM = log.Fields(m)
}

// Debugf log with debug level
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.preM != nil {
		l.Logger.WithFields(l.preM).Debugf(format, v...)
	} else {
		l.Logger.Debugf(format, v...)
	}
}

// Infof log with info level
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.preM != nil {
		l.Logger.WithFields(l.preM).Infof(format, v...)
	} else {
		l.Logger.Infof(format, v...)
	}
}

// Warnf log with warn level
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.preM != nil {
		l.Logger.WithFields(l.preM).Warnf(format, v...)
		if l.welog != nil {
			l.welog.WithFields(l.preM).Warnf(format, v...)
		}
	} else {
		l.Logger.Warnf(format, v...)
		if l.welog != nil {
			l.welog.Warnf(format, v...)
		}
	}
}

// Errorf log with error level
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.preM != nil {
		l.Logger.WithFields(l.preM).Errorf(format, v...)
		if l.welog != nil {
			l.welog.WithFields(l.preM).Errorf(format, v...)
		}
	} else {
		l.Logger.Errorf(format, v...)
		if l.welog != nil {
			l.welog.Errorf(format, v...)
		}
	}
}

// Fatalf log with fatal level
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.preM != nil {
		l.Logger.WithFields(l.preM).Fatalf(format, v...)
		if l.welog != nil {
			l.welog.WithFields(l.preM).Fatalf(format, v...)
		}
	} else {
		l.Logger.Fatalf(format, v...)
		if l.welog != nil {
			l.welog.Fatalf(format, v...)
		}
	}
}

// GetEntryWithFields return an entry with fileds, inherit origin fields
func (l *Logger) GetEntryWithFields(m map[string]interface{}) *Entry {
	entry := &Entry{Entry: l.getEntry(m, false)}
	if l.welog != nil {
		entry.weEntry = l.getEntry(m, true)
	}

	return entry
}

func (l *Logger) getEntry(m map[string]interface{}, isWE bool) *log.Entry {
	if m == nil {
		m = make(map[string]interface{})
	}

	if l.preM != nil {
		for k, v := range l.preM {
			if _, ok := m[k]; !ok {
				// 对preM中定义且m中没有的元素赋值
				m[k] = v
			}
		}
	}

	if isWE {
		return l.welog.WithFields(log.Fields(m))
	}

	return l.Logger.WithFields(log.Fields(m))
}

// GetField return field of key
func (l *Logger) GetField(key string) interface{} {
	return l.preM[key]
}

// Init initilize stdLog
func Init(level string, out io.Writer, flag int, relativePath string, errout io.Writer) {
	stdLog = NewWithRelativePath(level, out, flag, relativePath, errout)
}

// GetStdLog export stdLog
func GetStdLog() *Logger {
	return stdLog
}

// Predefine stdLog with fixed extra record
func Predefine(m map[string]interface{}) {
	stdLog.mu.Lock()
	defer stdLog.mu.Unlock()
	stdLog.preM = log.Fields(m)
}

// Debugf stdLog with debug level
func Debugf(format string, v ...interface{}) {
	if stdLog.preM != nil {
		stdLog.Logger.WithFields(stdLog.preM).Debugf(format, v...)
	} else {
		stdLog.Logger.Debugf(format, v...)
	}
}

// Infof stdLog with info level
func Infof(format string, v ...interface{}) {
	if stdLog.preM != nil {
		stdLog.Logger.WithFields(stdLog.preM).Infof(format, v...)
	} else {
		stdLog.Logger.Infof(format, v...)
	}
}

// Warnf stdLog with warn level
func Warnf(format string, v ...interface{}) {
	if stdLog.preM != nil {
		stdLog.Logger.WithFields(stdLog.preM).Warnf(format, v...)
		if stdLog.welog != nil {
			stdLog.welog.WithFields(stdLog.preM).Warnf(format, v...)
		}
	} else {
		stdLog.Logger.Warnf(format, v...)
		if stdLog.welog != nil {
			stdLog.welog.Warnf(format, v...)
		}
	}
}

// Errorf stdLog with error level
func Errorf(format string, v ...interface{}) {
	if stdLog.preM != nil {
		stdLog.Logger.WithFields(stdLog.preM).Errorf(format, v...)
		if stdLog.welog != nil {
			stdLog.welog.WithFields(stdLog.preM).Errorf(format, v...)
		}
	} else {
		stdLog.Logger.Errorf(format, v...)
		if stdLog.welog != nil {
			stdLog.welog.Errorf(format, v...)
		}
	}
}

// Fatalf stdLog with fatal level
func Fatalf(format string, v ...interface{}) {
	if stdLog.preM != nil {
		stdLog.Logger.WithFields(stdLog.preM).Fatalf(format, v...)
		if stdLog.welog != nil {
			stdLog.welog.WithFields(stdLog.preM).Fatalf(format, v...)
		}
	} else {
		stdLog.Logger.Fatalf(format, v...)
		if stdLog.welog != nil {
			stdLog.welog.Fatalf(format, v...)
		}
	}
}

// GetEntryWithFields return an entry with fileds, inherit origin fields from stdLog
func GetEntryWithFields(m map[string]interface{}) *Entry {
	entry := &Entry{Entry: stdLog.getEntry(m, false)}
	if stdLog.welog != nil {
		entry.weEntry = stdLog.getEntry(m, true)
	}

	return entry
}

// GetField return field of key from stdLog
func GetField(key string) interface{} {
	return stdLog.preM[key]
}
