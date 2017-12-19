package mylog

import (
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FileLineHook hook to get fileline
type FileLineHook struct {
	IsRelative   bool
	IgnorePrefix string
}

// Levels get levels
func (hook FileLineHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire execute hook
func (hook FileLineHook) Fire(entry *log.Entry) error {
	// 指针数组长度决定能存多少层堆栈
	pc := make([]uintptr, 3, 3)
	// 跳过最底下几层，至少7层才能跳到当前程序目录
	skip := 7
	// 返回堆栈调用层级
	cnt := runtime.Callers(skip, pc)

	// logrus的不同日志打印方法调用的堆栈层级不同
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		// 当调用的不是logrus的库文件或pkg/mylog时，打印当前文件信息
		if !strings.Contains(name, "github.com/sirupsen/logrus") && !strings.Contains(name, "util/mylog") {
			file, line := fu.FileLine(pc[i] - 1)
			if hook.IsRelative {
				// 使用相对路径
				entry.Data["file"] = strings.TrimPrefix(file, hook.IgnorePrefix)
			} else {
				// 使用文件名
				entry.Data["file"] = path.Base(file)
			}

			entry.Data["line"] = line
			break
		}
	}

	return nil
}
