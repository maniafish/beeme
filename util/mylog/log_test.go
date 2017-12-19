package mylog

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLog(t *testing.T) {
	log := New("DEBUG", os.Stderr, 0)
	log.Debugf("test log debug")
	log.Infof("test log info")
	log.Warnf("test log warn")
	log.Errorf("test log error")
}

func TestGetEntryWithFields(t *testing.T) {
	log := New("DEBUG", os.Stderr, Lfile)
	m := map[string]interface{}{
		"sdkid": "demosdk",
	}

	log.Predefine(m)
	log.Infof("test log")
	entry := log.GetEntryWithFields(nil)
	entry.Infof("test log with nil map")
	entry1 := entry.GetEntryWithFields(nil)
	entry1.Infof("test entry with nil map")
	entry2 := entry.GetEntryWithFields(map[string]interface{}{
		"entry2": "yes",
	})

	entry1.Infof("test entry1 with map")
	entry2.Infof("test entry2 with map")
}

func TestGetField(t *testing.T) {
	log := New("DEBUG", os.Stderr, 0)
	t.Logf("log field: %v", log.GetField("test"))
	entry := log.GetEntryWithFields(nil)
	t.Logf("entry field: %v", entry.GetField("test"))
}

func TestStdLog(t *testing.T) {
	_, lastFilePath, _, _ := runtime.Caller(1)
	ignoredDirPrefix := filepath.Dir(filepath.Dir(lastFilePath))
	Init("DEBUG", os.Stderr, Lfile|Lrelative|Ljson, ignoredDirPrefix, nil)
	Debugf("new debug")
	Infof("new info")
	Warnf("new warning")
	Errorf("new error")

	m := map[string]interface{}{
		"pid":  os.Getpid(),
		"mark": "stdlog",
	}

	Predefine(m)
	entry := GetEntryWithFields(nil)
	entry.Debugf("test entry debug")
	entry.Infof("test entry info")
	entry.Warnf("test entry warn")
	entry.Errorf("test entry error")
}

func TestWElog(t *testing.T) {
	_, lastFilePath, _, _ := runtime.Caller(1)
	ignoredDirPrefix := filepath.Dir(filepath.Dir(lastFilePath))
	Init("DEBUG", os.Stderr, Lfile|Lrelative|Ljson|Lwelog, ignoredDirPrefix, os.Stdout)
	Debugf("test stdlog no-prem debug")
	Infof("test stdlog no-prem info")
	Warnf("test stdlog no-prem warn")
	Errorf("test stdlog no-prem error")

	m := map[string]interface{}{
		"pid":  os.Getpid(),
		"mark": "welog",
	}

	Predefine(m)
	log := GetStdLog()
	log.Debugf("test debug")
	log.Infof("test info")
	log.Warnf("test warn")
	log.Errorf("test error")

	Debugf("test stdlog debug")
	Infof("test stdlog info")
	Warnf("test stdlog warn")
	Errorf("test stdlog error")

	entry := log.GetEntryWithFields(map[string]interface{}{
		"entry": "yes",
	})

	entry.Debugf("test entry debug")
	entry.Infof("test entry info")
	entry.Warnf("test entry warn")
	entry.Errorf("test entry error")

	entry2 := entry.GetEntryWithFields(map[string]interface{}{
		"entry2": "no",
	})

	entry2.Debugf("test entry debug")
	entry2.Infof("test entry info")
	entry2.Warnf("test entry warn")
	entry2.Errorf("test entry error")
}
