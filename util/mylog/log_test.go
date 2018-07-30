package mylog

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func checkBuffer(get *bytes.Buffer, expect map[string]string) {
	m := make(map[string]string)
	for _, v := range strings.Split(get.String(), " ") {
		if !strings.Contains(v, "=") {
			continue
		}

		kvArr := strings.Split(v, "=")
		key := strings.TrimSpace(kvArr[0])
		val := strings.TrimSpace(kvArr[1])
		if kvArr[1][0] == '"' {
			val, _ = strconv.Unquote(val)
		}

		m[key] = val
	}

	for k, v := range expect {
		// time 允许正负1s误差
		if k == "time" {
			t, err := time.ParseInLocation("2006-01-02T15:04:05", m[k][:19], time.Local)
			So(err, ShouldBeNil)
			e, err := time.ParseInLocation("2006-01-02T15:04:05", v, time.Local)
			So(err, ShouldBeNil)
			So(t, ShouldHappenBetween, e.Add(-1*time.Second), e.Add(1*time.Second))
		} else {
			So(m[k], ShouldEqual, v)
		}
	}
}

func checkJSONBuffer(get *bytes.Buffer, expect map[string]string) {
	m := make(map[string]string)
	err := json.Unmarshal(get.Bytes(), &m)
	So(err, ShouldBeNil)
	for k, v := range expect {
		if k == "time" {
			t, err := time.ParseInLocation("2006-01-02T15:04:05", m[k][:19], time.Local)
			So(err, ShouldBeNil)
			e, err := time.ParseInLocation("2006-01-02T15:04:05", v, time.Local)
			So(err, ShouldBeNil)
			So(t, ShouldHappenBetween, e.Add(-1*time.Second), e.Add(1*time.Second))
		} else {
			So(m[k], ShouldEqual, v)
		}
	}
}

func TestLog(t *testing.T) {
	var buffer bytes.Buffer
	Init("DEBUG", &buffer, 0, "", nil)
	Predefine(map[string]interface{}{
		"tag": "stdlog",
	})

	Convey("Subject GetField\n", t, func() {
		So(GetField("tag"), ShouldEqual, "stdlog")
	})

	Convey("Subject: test debug stdlog\n", t, func() {
		buffer.Reset()
		Debugf("debug1")
		Println(buffer.String())
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "debug1",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test info stdlog\n", t, func() {
		buffer.Reset()
		Infof("info1")
		checkBuffer(&buffer, map[string]string{
			"level": "info",
			"msg":   "info1",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test warning stdlog\n", t, func() {
		buffer.Reset()
		Warnf("warning1")
		checkBuffer(&buffer, map[string]string{
			"level": "warning",
			"msg":   "warning1",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test error stdlog\n", t, func() {
		buffer.Reset()
		Errorf("error1")
		checkBuffer(&buffer, map[string]string{
			"level": "error",
			"msg":   "error1",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	log := GetStdLog()
	Convey("Subject: test debug log\n", t, func() {
		buffer.Reset()
		log.Debugf("debug2")
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "debug2",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test info log\n", t, func() {
		buffer.Reset()
		log.Infof("info2")
		checkBuffer(&buffer, map[string]string{
			"level": "info",
			"msg":   "info2",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test warning log\n", t, func() {
		buffer.Reset()
		log.Warnf("warning2")
		checkBuffer(&buffer, map[string]string{
			"level": "warning",
			"msg":   "warning2",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test error log\n", t, func() {
		buffer.Reset()
		log.Errorf("error2")
		checkBuffer(&buffer, map[string]string{
			"level": "error",
			"msg":   "error2",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"tag":   "stdlog",
		})
	})

	Convey("Subject: test invalid loglevel\n", t, func() {
		log := New("invalid", &buffer, 0)
		buffer.Reset()
		log.Debugf("default_debug")
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "default_debug",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})
	})
}

func TestLogFlag(t *testing.T) {
	var buffer bytes.Buffer
	// fileline hook 会自动跳过mylog包
	/*
		Convey("Subject: test flag Lfile\n", t, func() {
			log := New("DEBUG", &buffer, Lfile)
			buffer.Reset()
			log.Debugf("file")
			checkBuffer(&buffer, map[string]string{
				"level": "debug",
				"msg":   "file",
				"time":  time.Now().Format("2006-01-02T15:04:05"),
				"file":  "log_test.go",
				"line":  "184",
			})
		})

		Convey("Subject: test flag Lrelative\n", t, func() {
			_, thisFilePath, _, _ := runtime.Caller(0)
			ignoredDirPrefix := filepath.Dir(filepath.Dir(thisFilePath))
			log := NewWithRelativePath("DEBUG", &buffer, Lfile|Lrelative, ignoredDirPrefix, nil)
			buffer.Reset()
			log.Debugf("relative")
			checkBuffer(&buffer, map[string]string{
				"level": "debug",
				"msg":   "relative",
				"time":  time.Now().Format("2006-01-02T15:04:05"),
				"file":  "/mylog/log_test.go",
				"line":  "199",
			})
		})
	*/

	Convey("Subject: test flag Ljson\n", t, func() {
		log := New("DEBUG", &buffer, Ljson)
		buffer.Reset()
		log.Debugf("json")
		checkJSONBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "json",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})

		buffer.Reset()
		log.Infof("json_info")
		checkJSONBuffer(&buffer, map[string]string{
			"level": "info",
			"msg":   "json_info",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})
	})

	Convey("Subject: test flag Lwelog\n", t, func() {
		var weBuffer bytes.Buffer
		log := NewWithRelativePath("DEBUG", &buffer, Lwelog, "", &weBuffer)
		buffer.Reset()
		log.Debugf("welog")
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "welog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})

		checkBuffer(&weBuffer, map[string]string{
			"level": "",
			"msg":   "",
		})

		buffer.Reset()
		weBuffer.Reset()
		log.Warnf("warninglog")
		checkBuffer(&buffer, map[string]string{
			"level": "warning",
			"msg":   "warninglog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})

		checkBuffer(&weBuffer, map[string]string{
			"level": "warning",
			"msg":   "warninglog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})

		buffer.Reset()
		weBuffer.Reset()
		log.Errorf("errorlog")
		checkBuffer(&buffer, map[string]string{
			"level": "error",
			"msg":   "errorlog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})

		checkBuffer(&weBuffer, map[string]string{
			"level": "error",
			"msg":   "errorlog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})
	})

	Convey("Subject: test flag Lwelog with prem\n", t, func() {
		var weBuffer bytes.Buffer
		log := NewWithRelativePath("DEBUG", &buffer, Lwelog, "", &weBuffer)
		log.Predefine(map[string]interface{}{
			"id": 110,
		})

		buffer.Reset()
		log.Debugf("welog")
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "welog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"id":    "110",
		})

		checkBuffer(&weBuffer, map[string]string{
			"level": "",
			"msg":   "",
			"id":    "",
		})

		buffer.Reset()
		weBuffer.Reset()
		log.Warnf("warninglog")
		checkBuffer(&buffer, map[string]string{
			"level": "warning",
			"msg":   "warninglog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"id":    "110",
		})

		checkBuffer(&weBuffer, map[string]string{
			"level": "warning",
			"msg":   "warninglog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"id":    "110",
		})

		buffer.Reset()
		weBuffer.Reset()
		log.Errorf("errorlog")
		checkBuffer(&buffer, map[string]string{
			"level": "error",
			"msg":   "errorlog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"id":    "110",
		})

		checkBuffer(&weBuffer, map[string]string{
			"level": "error",
			"msg":   "errorlog",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"id":    "110",
		})
	})
}

func TestGetEntryWithFields(t *testing.T) {
	var buffer bytes.Buffer
	Init("DEBUG", &buffer, 0, "", nil)
	Convey("Subject: test GetEntryWithFields\n", t, func() {
		nlog := GetEntryWithFields(map[string]interface{}{
			"id": 123,
		})

		buffer.Reset()
		nlog.Debugf("get_entry_with_fields")
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "get_entry_with_fields",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
			"id":    "123",
		})

		So(nlog.GetField("id").(int), ShouldEqual, 123)
		n2log := nlog.GetEntryWithFields(map[string]interface{}{
			"id":     456,
			"new_id": 457,
		})

		buffer.Reset()
		n2log.Debugf("get_entry_with_fields2")
		checkBuffer(&buffer, map[string]string{
			"level":  "debug",
			"msg":    "get_entry_with_fields2",
			"time":   time.Now().Format("2006-01-02T15:04:05"),
			"id":     "456",
			"new_id": "457",
		})
	})
}

func TestGetCtxEntry(t *testing.T) {
	var buffer bytes.Buffer
	log := New("DEBUG", &buffer, 0)
	Convey("Subject: test GetCtxEntry\n", t, func() {
		nctx := SetCtxEntry(context.Background(), log.GetEntryWithFields(map[string]interface{}{
			"requestid": "ctx",
		}))

		nlog := GetCtxEntry(nctx, log)
		buffer.Reset()
		nlog.Debugf("get_ctx_entry")
		checkBuffer(&buffer, map[string]string{
			"level":     "debug",
			"msg":       "get_ctx_entry",
			"time":      time.Now().Format("2006-01-02T15:04:05"),
			"requestid": "ctx",
		})
	})

	Convey("Subject: test nil GetCtxEntry\n", t, func() {
		nlog := GetCtxEntry(context.Background(), log)
		buffer.Reset()
		nlog.Debugf("get_nil_ctx_entry")
		checkBuffer(&buffer, map[string]string{
			"level": "debug",
			"msg":   "get_nil_ctx_entry",
			"time":  time.Now().Format("2006-01-02T15:04:05"),
		})
	})
}
