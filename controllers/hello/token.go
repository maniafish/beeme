package hello

import (
	"beeme/controllers"
	"beeme/util/sha1str"
	"beeme/util/sort"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Controller Operations about hello
type Controller struct {
	controllers.Controller
}

// TokenRequest request for token
type TokenRequest struct {
	Sign      string `sort:"-"`
	Timestamp string `sort:"timestamp"`
	Nonce     string `sort:"nonce"`
	Echostr   string `sort:"-"`
	Token     string `sort:"token"`
}

func sha1Sign(v interface{}) string {
	values := sort.GetValues(v)
	strToSign := values.Encode("")
	return sha1str.HexString(strToSign)
}

// Get get
// @Title Get
// @Description verify js token, doc site: https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319
// @router / [get]
func (c *Controller) Get() {
	logPrefix := "hello.Get()"
	req := TokenRequest{
		Sign:      c.GetString("signature"),
		Timestamp: c.GetString("timestamp"),
		Nonce:     c.GetString("nonce"),
		Echostr:   c.GetString("echostr"),
		Token:     beego.AppConfig.String("apps::AppID"),
	}

	logs.Info("%s.TokenRequest: %+v", logPrefix, req)
	expect := sha1Sign(req)
	if expect != req.Sign {
		logs.Info("%s.verifySign: expect sign: %s", logPrefix, expect)
		c.Ctx.WriteString("invalid sign")
		return
	}

	c.Ctx.WriteString(req.Echostr)
	return
}
