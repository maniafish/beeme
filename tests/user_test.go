package test

import (
	"beeme/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func testGetUser(t *testing.T, id int, expCode int, expect interface{}) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("/v1/user/%d", id), nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test User Get\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, expCode)
		})

		if expCode == 200 {
			Convey("Should Have Valid Body\n", func() {
				user := &models.User{}
				ret, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				err = json.Unmarshal(ret, user)
				So(err, ShouldBeNil)
				So(user, ShouldResemble, expect)
			})
		}
	})
}

// TestUser is a sample to run an user api test
func TestUser(t *testing.T) {
	type UID struct {
		Data int `json:"uid"`
	}

	uid := &UID{}
	user := &models.User{
		Username: "test",
		Password: "pass",
		Gender:   "male",
		Age:      1,
		Address:  "netease",
		Email:    "mania",
	}

	// Test Post
	body, _ := json.Marshal(user)
	r, _ := http.NewRequest("POST", "/v1/user/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test User Post\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should Have Valid Body\n", func() {
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			err = json.Unmarshal(ret, uid)
			So(err, ShouldBeNil)
			So(uid.Data, ShouldBeGreaterThan, 0)
		})
	})

	// TestGet
	user.ID = uid.Data
	testGetUser(t, user.ID, 200, user)

	// TestPut
	user.Username += "+1"
	body, _ = json.Marshal(user)
	r, _ = http.NewRequest("PUT", fmt.Sprintf("/v1/user/%d", user.ID), bytes.NewReader(body))
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test User Update\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should Have Valid Body\n", func() {
			userGet := &models.User{}
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			err = json.Unmarshal(ret, userGet)
			So(err, ShouldBeNil)
			So(userGet, ShouldResemble, user)
		})

	})

	testGetUser(t, user.ID, 200, user)

	// TestDelete
	r, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/user/%d", user.ID), nil)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test User Update\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should Have Valid Body\n", func() {
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			err = json.Unmarshal(ret, uid)
			So(err, ShouldBeNil)
			So(uid.Data, ShouldEqual, user.ID)
		})
	})

	testGetUser(t, user.ID, 404, nil)
}
