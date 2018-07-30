package controllers

import (
	"beeme/util/counter"
	"beeme/util/mylog"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/astaxie/beego"
)

var reqid = counter.New()

// Controller base controller implement log wrap, request and response
type Controller struct {
	beego.Controller
	Logger mylog.MyLogger
	Name   string
}

// Prepare init Controller.Logger(set requestid)
func (c *Controller) Prepare() {
	c.Logger = mylog.GetEntryWithFields(map[string]interface{}{
		"requestid": fmt.Sprintf("%s.%v", c.Name, reqid.Incr()),
	})
}

// ServeJSON return json-data
func (c *Controller) ServeJSON(code int, v interface{}) {
	c.Ctx.Output.Status = code
	c.Data["json"] = v
	c.Controller.ServeJSON()
	c.ServerLog("json")
}

// ServeXML return xml-data
func (c *Controller) ServeXML(code int, v interface{}) {
	c.Ctx.Output.Status = code
	c.Data["xml"] = v
	c.Controller.ServeXML()
	c.ServerLog("xml")
}

// ServeString return string-data
func (c *Controller) ServeString(code int, v string) {
	c.Ctx.Output.Status = code
	c.Data["string"] = v
	c.Controller.Ctx.WriteString(v)
	c.ServerLog("string")
}

// ServerLog logf request and response
func (c *Controller) ServerLog(retType string) {
	var resp string
	switch retType {
	case "xml":
		b, e := xml.Marshal(c.Data["xml"])
		if e != nil {
			c.Logger.Errorf("return err: %v", e)
		} else {
			resp = string(b)
		}
	case "json":
		b, e := json.Marshal(c.Data["json"])
		if e != nil {
			c.Logger.Errorf("return err: %v", e)
		} else {
			resp = string(b)
		}
	case "string":
		resp, _ = c.Data["string"].(string)
	default:
		c.Logger.Errorf("invalid retType: %v", retType)
	}

	c.Logger.Infof("header: %v, method: %v, url: %v, body: %s, resp: %v", c.Ctx.Request.Method, c.Ctx.Request.Method, c.Ctx.Request.RequestURI, c.Ctx.Input.RequestBody, resp)
}
