package robot

import (
	"beeme/conf"
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
	keyID := robotid.Incr() % uint64(len(conf.Config.TulingKeys))
	return conf.Config.TulingKeys[keyID]
}

// Get get
// @Title Get
// @Description get question
// @Param	ask		path 	string	true		"The question you ask"
// @Success 200 {object} robot.Resp
// @router /ask/:ask [get]
func (b *Controller) Get() {
	question := b.GetString(":ask")
	req, _ := json.Marshal(&models.RobotReq{
		Key:    getKey(),
		Info:   question,
		UserID: "123",
	})

	resp, err := httplib.Post(conf.Config.TulingURL).Body(req).Bytes()
	if err != nil {
		b.Response(501, "external error")
	}

	respObj := &models.RobotResp{}
	err = json.Unmarshal(resp, respObj)
	if err != nil {
		b.Response(501, "external error")
	}

	if !respObj.IsValid() {
		logs.Warning("respObj: %+v", respObj)
		b.Response(400, fmt.Sprintf("%v, %v", respObj.Code, respObj.Text))
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
