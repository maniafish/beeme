package test

import (
	_ "beeme/routers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUser is a sample to run an user api test
func TestUser(t *testing.T) {
	type UID struct {
		Data int `json:"uid"`
	}
	r, _ := http.NewRequest("GET", fmt.Sprintf("/v1/user/get/test"), nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Subject: Test User Invalid\n", t, func() {
		Convey("Status Code Should Be 403\n", func() {
			So(w.Code, ShouldEqual, 403)
		})
	})

	r, _ = http.NewRequest("GET", "/v1/user/get/123", nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Subject: Test GetUser\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("The Result Should Be 123\n", func() {
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			uid := &UID{}
			err = json.Unmarshal(ret, uid)
			So(err, ShouldBeNil)
			So(uid.Data, ShouldEqual, 123)
		})
	})
}
