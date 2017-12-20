package test

import (
	"beeme/controllers/wechat"
	"beeme/models"
	"beeme/util/xmls"
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

// TestMessage is a sample to run wechat_message test
func TestMessage(t *testing.T) {
	req := &wechat.MessageReq{
		ToUserName:   xmls.CharData("to"),
		FromUserName: xmls.CharData("from"),
		CreateTime:   time.Now().Unix(),
		MsgType:      xmls.CharData("teXt"),
		MsgID:        123456,
		Content:      xmls.CharData("小O"),
	}

	body, _ := xml.Marshal(req)
	r, _ := http.NewRequest("POST", "/v1/js", bytes.NewReader(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	Convey("Subject: Test Message Post\n", t, func() {
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should Have Valid Body\n", func() {
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			resp := &wechat.MessageResp{}
			err = xml.Unmarshal(ret, resp)
			So(err, ShouldBeNil)
			So(string(resp.Content), ShouldEqual, models.DefaultErrMsg)
		})
	})
}
