package demo

import (
	"beeme/controllers"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/astaxie/beego"
)

// FileController Operations about Files
type FileController struct {
	controllers.Controller
}

// FileResponse response struct
type FileResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Prepare set controller param
func (u *FileController) Prepare() {
	u.Name = "file_operator"
	u.Controller.Prepare()
}

// PostFlow postflow
// @Title PostFlow
// @Description 处理流式上传文件，存入本地
// @router /flow [post]
func (u *FileController) PostFlow() {
	tofile := path.Join(beego.AppConfig.String("apps::FileDir"), fmt.Sprintf("%v", time.Now().Unix()))
	tf, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		u.Logger.Infof("OpenFile err: %v", err)
		u.ServeJSON(200, &FileResponse{-1, "internal error"})
		return
	}

	defer tf.Close()
	_, err = tf.Write(u.Ctx.Input.RequestBody)
	if err != nil {
		u.Logger.Infof("WriteFile err: %v", err)
		u.ServeJSON(200, &FileResponse{-1, "internal error"})
		return
	}

	u.ServeJSON(200, &FileResponse{0, tofile})
	return
}
