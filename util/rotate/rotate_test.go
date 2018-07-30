package rotate

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRotateTime(t *testing.T) {
	t1, _ := time.Parse("2006-01-02 15:04:05", "2017-02-28 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-02-28 00:00:00")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 00:00:00")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-31 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-31 00:00:00")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-30 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-30 00:00:00")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 00:00:00")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-04-30 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-04-30 00:00:00")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-12-31 23:59:59")
	testGetRotateTime(t, t1)
	t1, _ = time.Parse("2006-01-02 15:04:05", "2017-12-31 00:00:00")
	testGetRotateTime(t, t1)
	t1 = time.Now()
	testGetRotateTime(t, t1)
}

func testGetRotateTime(t *testing.T, t1 time.Time) {
	prefix := fmt.Sprintf("%v: ", t1)
	Convey(prefix+"Subject: Test Daily Rotate\n", t, func() {
		dailyTime := getRotateTime(Daily, t1)
		timePrefix := fmt.Sprintf("%v: ", dailyTime)

		Convey(timePrefix+"Date Should Equal Now\n", func() {
			So(dailyTime.Year(), ShouldEqual, t1.Year())
			So(dailyTime.Month(), ShouldEqual, t1.Month())
			So(dailyTime.Day(), ShouldEqual, t1.Day())
		})

		Convey(timePrefix+"Time Should be 23:59:59\n", func() {
			So(dailyTime.Hour(), ShouldEqual, 23)
			So(dailyTime.Minute(), ShouldEqual, 59)
			So(dailyTime.Second(), ShouldEqual, 59)
		})
	})

	Convey(prefix+"Subject: Test Weekly Rotate\n", t, func() {
		weeklyTime := getRotateTime(Weekly, t1)
		timePrefix := fmt.Sprintf("%v: ", weeklyTime)

		Convey(timePrefix+"Week Should Equal Sunday\n", func() {
			So(weeklyTime.Weekday(), ShouldEqual, time.Sunday)
		})

		Convey(timePrefix+"0 days <= weeklyTime - Now < 7 days\n", func() {
			duration := weeklyTime.Sub(t1)
			So(duration.Hours(), ShouldBeGreaterThanOrEqualTo, 0.0)
			So(duration.Hours(), ShouldBeLessThan, 7*24.0)
		})

		Convey(timePrefix+"Time Should be 23:59:59\n", func() {
			So(weeklyTime.Hour(), ShouldEqual, 23)
			So(weeklyTime.Minute(), ShouldEqual, 59)
			So(weeklyTime.Second(), ShouldEqual, 59)
		})
	})

	Convey(prefix+"Subject: Test Monthly Rotate\n", t, func() {
		monthlyTime := getRotateTime(Monthly, t1)
		timePrefix := fmt.Sprintf("%v: ", monthlyTime)

		Convey(timePrefix+"(Date + 1).Day Should Equal 1", func() {
			So(monthlyTime.AddDate(0, 0, 1).Day(), ShouldEqual, 1)
		})

		Convey(timePrefix+"Year Should Equal now.Year", func() {
			So(monthlyTime.Year(), ShouldEqual, t1.Year())
		})

		Convey(timePrefix+"Month Should Equal now.Month", func() {
			So(monthlyTime.Month(), ShouldEqual, t1.Month())
		})

		Convey(timePrefix+"Time Should be 23:59:59\n", func() {
			So(monthlyTime.Hour(), ShouldEqual, 23)
			So(monthlyTime.Minute(), ShouldEqual, 59)
			So(monthlyTime.Second(), ShouldEqual, 59)
		})
	})

	// as same as daily rotate
	Convey(prefix+"Subject: Test Default Rotate\n", t, func() {
		defaultTime := getRotateTime(Daily, t1)
		timePrefix := fmt.Sprintf("%v: ", defaultTime)

		Convey(timePrefix+"Date Should Equal Now\n", func() {
			So(defaultTime.Year(), ShouldEqual, t1.Year())
			So(defaultTime.Month(), ShouldEqual, t1.Month())
			So(defaultTime.Day(), ShouldEqual, t1.Day())
		})

		Convey(timePrefix+"Time Should be 23:59:59\n", func() {
			So(defaultTime.Hour(), ShouldEqual, 23)
			So(defaultTime.Minute(), ShouldEqual, 59)
			So(defaultTime.Second(), ShouldEqual, 59)
		})
	})
}
