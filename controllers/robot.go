package controllers

import (
	"beeme/conf"
	"beeme/models"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

// BotController Operations about BotChat
type BotController struct {
	beego.Controller
}

// BotResp response
type BotResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Get get
// @Title Get
// @Description get question
// @Param	ask		path 	string	true		"The question you ask"
// @Success 200 {object} controllers.BotResp
// @router /ask/:ask [get]
func (b *BotController) Get() {
	question := b.GetString(":ask")
	req, _ := json.Marshal(&models.RobotReq{
		Key:    conf.Config.TulingKeys[0],
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
func (b *BotController) Response(code int, msg string) {
	resp := &BotResp{
		Code: code,
		Msg:  msg,
	}

	b.Data["json"] = resp
	b.ServeJSON()
}
