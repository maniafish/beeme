package hello

import (
	"beeme/util/sha1str"
	"beeme/util/sort"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Controller Operations about hello
type Controller struct {
	beego.Controller
}

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
// @Description verify js token
// @router / [get]
func (c *Controller) Get() {
	logPrefix := "hello.Get()"
	req := TokenRequest{
		Sign:      c.GetString("signature"),
		Timestamp: c.GetString("timestamp"),
		Nonce:     c.GetString("nonce"),
		Echostr:   c.GetString("echostr"),
		Token:     c.GetString("token"),
	}

	logs.Info("%s.TokenRequest: %+v", logPrefix, req)
	if req.Token != beego.AppConfig.String("apps::AppID") {
		c.Ctx.WriteString("invalid token")
		return
	}

	expect := sha1Sign(req)
	if expect != req.Sign {
		logs.Info("%s.verifySign: expect sign: %s", logPrefix, expect)
		c.Ctx.WriteString("invalid sign")
		return
	}

	c.Ctx.WriteString(req.Echostr)
	return
}
