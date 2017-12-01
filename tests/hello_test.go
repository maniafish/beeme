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
	v.Set("signature", "4223754551c082176b29c847ef79bb297378f2f7")
	v.Set("timestamp", "1512121658")
	v.Set("nonce", "no")
	v.Set("echostr", "success")
	v.Set("token", beego.AppConfig.String("apps::AppID"))
	url := "/v1" + "?" + v.Encode()
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
