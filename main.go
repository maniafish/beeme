package main

import (
	"beeme/models"
	_ "beeme/routers"
	"beeme/util/mylog"

	"github.com/astaxie/beego"
)

var version = "version: unknown"

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	mylog.Infof(version)
	models.Init()
	beego.Run()
}
