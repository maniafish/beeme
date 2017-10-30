package routers

// @APIVersion 1.0.0 // @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

import (
	"beeme/controllers"
	"beeme/controllers/robot"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/robot",
			beego.NSInclude(
				&robot.Controller{},
			),
		),
	)
	beego.AddNamespace(ns)
}
