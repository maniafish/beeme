package robot

import (
	"beeme/models"
	"beeme/util/counter"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

// Controller Operations about BotChat
type Controller struct {
	beego.Controller
}

// Resp response
type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var robotid = counter.New()

func getKey() string {
	keyID := robotid.Incr() % uint64(len(beego.AppConfig.Strings("apps::TulingKeys")))
	return beego.AppConfig.Strings("apps::TulingKeys")[keyID]
}

// Get get
// @Title Get
// @Description get question
// @Param	ask		path 	string	true		"The question you ask"
// @Success 200 {object} robot.Resp
// @router /ask/:ask [get]
func (b *Controller) Get() {
	logPrefix := "robot.Get()"
	question := b.GetString(":ask")
	req, _ := json.Marshal(&models.RobotReq{
		Key:    getKey(),
		Info:   question,
		UserID: "",
	})

	resp, err := httplib.Post(beego.AppConfig.String("apps::TulingURL")).Body(req).Bytes()
	if err != nil {
		logs.Error("%s.Post err: %+v", logPrefix, err)
		b.Response(501, "external error")
		return
	}

	respObj := &models.RobotResp{}
	err = json.Unmarshal(resp, respObj)
	if err != nil {
		logs.Error("%s.PostResp err: %+v", logPrefix, err)
		b.Response(501, "external error")
		return
	}

	if !respObj.IsValid() {
		logs.Warning("%s.respObj Invalid: %+v", logPrefix, respObj)
		b.Response(400, fmt.Sprintf("%v, %v", respObj.Code, respObj.Text))
		return
	}

	b.Response(200, respObj.Text)
}

// Response return
func (b *Controller) Response(code int, msg string) {
	resp := &Resp{
		Code: code,
		Msg:  msg,
	}

	b.Data["json"] = resp
	b.ServeJSON()
}
