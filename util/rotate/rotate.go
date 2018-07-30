package rotate

import (
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"beeme/util/bg"
	"beeme/util/mylog"
)

const (
	// Daily do logrotate everyday at 23:59:59
	Daily = iota
	// Weekly do logrotate every week on Sunday at 23:59:59
	Weekly
	// Monthly do logrotate every month on last day of month at 23:59:59
	Monthly
)

func getRotateTime(policy int, t time.Time) time.Time {
	switch policy {
	case Daily:
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	case Weekly:
		lastDayWeekly := t.AddDate(0, 0, (7-int(t.Weekday()))%7)
		return time.Date(lastDayWeekly.Year(), lastDayWeekly.Month(), lastDayWeekly.Day(), 23, 59, 59, 0, lastDayWeekly.Location())
	case Monthly:
		// first day next month - 1 day
		lastDayMonthly := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, -1)
		return time.Date(lastDayMonthly.Year(), lastDayMonthly.Month(), lastDayMonthly.Day(), 23, 59, 59, 0, lastDayMonthly.Location())
	default:
		// default policy is daily
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	}
}

// InitRotate (Mylogger: logger; name: log file name; policy: Daily, Weekly, Monthly, maxBackups: max log backups saved), reutrn a rotate logger
func InitRotate(l mylog.MyLogger, name string, policy, maxBackups int) *lumberjack.Logger {
	rotate := &lumberjack.Logger{
		Filename:   name,
		MaxSize:    1024 * 1024 * 1024, // make sure maxsize large enough
		MaxBackups: maxBackups,
		LocalTime:  true,
	}

	bg.Run(l, "rotate", func() {
		next := getRotateTime(policy, time.Now())
		if bg.SleepWithWaitCtx(next.Sub(time.Now())) {
			return
		}

		rotate.Rotate()
		time.Sleep(2 * time.Second) // make sure ending rotate after 00:00:00, in case of duplicate rotating
	})

	return rotate
}
