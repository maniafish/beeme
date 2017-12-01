package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

// TestToken is a sample to run token test
func TestToken(t *testing.T) {
	v := url.Values{}
	v.Set("signature", "c3d3ca6f7e92bf2928f5419d5458826c10ee5ed3")
	v.Set("timestamp", "1512121658")
	v.Set("nonce", "no")
	v.Set("echostr", "success")
	url := "/v1/js" + "?" + v.Encode()
	r, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("Subject: Test Token\n", t, func() {
		Convey("Status Code Should be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Data Should be echostr", func() {
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			So(string(ret), ShouldEqual, v.Get("echostr"))
		})
	})
}
