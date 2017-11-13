package test

import (
	"beeme/models"
	// for init routers
	_ "beeme/routers"
	"path"
	"path/filepath"
	"runtime"

	"github.com/astaxie/beego"
)

func init() {
	_, thisFilePath, _, _ := runtime.Caller(0)
	err := beego.LoadAppConfig(
		"ini",
		path.Join(filepath.Dir(filepath.Dir(thisFilePath)), "conf", "app.conf"),
	)

	if err != nil {
		panic(err)
	}

	models.Init()
}
