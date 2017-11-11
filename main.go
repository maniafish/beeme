package main

import (
	"beeme/models"
	_ "beeme/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var version = "version: unknown"

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	logs.Info(version)
	models.Init()
	beego.Run()
}
