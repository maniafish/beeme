package test

import (
	"beeme/models"
	_ "beeme/routers"
	"path"
	"path/filepath"
	"runtime"

	"github.com/astaxie/beego"
)

func init() {
	_, thisFilePath, _, _ := runtime.Caller(0)
	beego.LoadAppConfig(
		"ini",
		path.Join(filepath.Dir(filepath.Dir(thisFilePath)), "conf", "app.conf"))
	models.Init()
}
