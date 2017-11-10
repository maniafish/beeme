package test

import (
	_ "beeme/routers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

// TestRobot is a sample to run robot  api test
func TestRobot(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/robot/ask/test", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	fmt.Printf("Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Tuling\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})
	})
}
