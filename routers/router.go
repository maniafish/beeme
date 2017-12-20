package routers

// @APIVersion 1.0.0 // @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"beeme/controllers/robot"
	"beeme/controllers/user"
	"beeme/controllers/wechat"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/js",
			beego.NSInclude(
				&wechat.Controller{},
			),
		),
		beego.NSNamespace("/robot",
			beego.NSInclude(
				&robot.Controller{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&user.Controller{},
			),
		),
	)
	beego.AddNamespace(ns)
}
