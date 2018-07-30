package mylog

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// Logger implement MyLogger interface
type Logger struct {
	*log.Logger

	welog *log.Logger // logger for warning and error
	mu    sync.Mutex
	preM  log.Fields
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

// GetEntryWithFields return an entry(*Logger) with fileds, inherit origin fields
func (l *Logger) GetEntryWithFields(m map[string]interface{}) *Logger {
	entry := Logger{
		Logger: l.Logger,
		welog:  l.welog,
	}

	preM := make(map[string]interface{})
	for k, v := range l.preM {
		preM[k] = v
	}

	for k, v := range m {
		preM[k] = v
	}

	entry.Predefine(preM)
	return &entry
}

// GetField return field of key
func (l *Logger) GetField(key string) interface{} {
	return l.preM[key]
}
