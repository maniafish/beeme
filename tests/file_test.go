package test

import (
	"beeme/controllers/demo"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

/*
func TestUploadFile(t *testing.T) {
	Convey("Subject: Test Upload File\n", t, func() {
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)
		fileWriter, err := bodyWriter.CreateFormFile("uploadfile", fmt.Sprintf("%v", time.Now().Unix()))
		So(err, ShouldBeNil)
		err = ioutil.WriteFile("/tmp/data", []byte("test\n"), 0666)
		So(err, ShouldBeNil)
		f, err := os.Open("/tmp/data")
		So(err, ShouldBeNil)
		defer f.Close()
		io.Copy(fileWriter, f)
		fmt.Printf("\n\n%s\n\n", bodyBuf)
		r, _ := http.NewRequest("POST", "/v1/uploadfile", bodyBuf)
		r.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should Have Valid Body\n", func() {
			resp := &demo.FileResponse{}
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			err = json.Unmarshal(ret, resp)
			So(err, ShouldBeNil)
			So(resp.Code, ShouldEqual, 0)
		})
	})
}
*/

func TestUploadFileFlow(t *testing.T) {
	Convey("Subject: Test Upload File Flow\n", t, func() {
		bodyBuf := &bytes.Buffer{}
		_, err := bodyBuf.Write([]byte("test12345678901234567890"))
		So(err, ShouldBeNil)
		r, _ := http.NewRequest("POST", "/v1/uploadfile/flow", bodyBuf)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		Convey("Status Code Should Be 200\n", func() {
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should Have Valid Body\n", func() {
			resp := &demo.FileResponse{}
			ret, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			err = json.Unmarshal(ret, resp)
			So(err, ShouldBeNil)
			So(resp.Code, ShouldEqual, 0)
			f, err := os.OpenFile(resp.Msg, os.O_RDONLY, 0666)
			So(err, ShouldBeNil)
			rb := make([]byte, 4)
			_, err = f.Read(rb)
			So(err, ShouldBeNil)
			So(string(rb), ShouldResemble, "test")
		})
	})
}
