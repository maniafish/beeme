package test

import (
	_ "beeme/routers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
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

	Convey(fmt.Sprintf("Subject: Test User Invalid\n"), t, func() {
		Convey("Status Code Should Be 404", func() {
			So(w.Code, ShouldEqual, 404)
		})
	})

	var wg sync.WaitGroup
	for i := 1; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r, _ := http.NewRequest("GET", fmt.Sprintf("/v1/user/get/%v", i), nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Convey(fmt.Sprintf("Subject: Test User %v\n", i), t, func() {
				Convey("Status Code Should Be 200", func() {
					So(w.Code, ShouldEqual, 200)
				})
				Convey(fmt.Sprintf("The Result Should Not Be %v", i), func() {
					ret, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					uid := &UID{}
					err = json.Unmarshal(ret, uid)
					So(err, ShouldBeNil)
					So(uid.Data, ShouldEqual, i)
				})
			})
		}(i)
	}

	wg.Wait()
}
