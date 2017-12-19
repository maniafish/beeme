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
	Logger  mylog.MyLogger
	retType string
	request interface{}
}

// Initialize init Controller.Logger(set requestid)
func (c *Controller) Initialize(prefix string) {
	c.Logger = mylog.GetEntryWithFields(map[string]interface{}{
		"requestid": fmt.Sprintf("%v.%v", prefix, reqid.Incr()),
	})
}

// SetRequest set request
func (c *Controller) SetRequest(v interface{}) {
	c.request = v
}

// GetRequest get request
func (c *Controller) GetRequest() interface{} {
	return c.request
}

// SetRetType set response type
func (c *Controller) SetRetType(v string) {
	c.retType = v
}

// GetRetType get response type
func (c *Controller) GetRetType() string {
	return c.retType
}

// ServeJSON return json-data
func (c *Controller) ServeJSON(v interface{}) {
	c.SetRetType("json")
	c.Data[c.retType] = v
	c.Controller.ServeJSON()
	c.ServerLog()
}

// ServeXML return xml-data
func (c *Controller) ServeXML(v interface{}) {
	c.SetRetType("xml")
	c.Data[c.retType] = v
	c.Controller.ServeXML()
	c.ServerLog()
}

// ServeString return string-data
func (c *Controller) ServeString(v string) {
	c.SetRetType("string")
	c.Data[c.retType] = v
	c.Controller.Ctx.WriteString(v)
	c.ServerLog()
}

// ServerLog logf request and response
func (c *Controller) ServerLog() {
	var resp string
	switch c.GetRetType() {
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
		c.Logger.Errorf("invalid retType: %v", c.retType)
	}

	c.Logger.Infof("header: %v, method: %v, url: %v, body: %s, resp: %v", c.Ctx.Request.Method, c.Ctx.Request.Method, c.Ctx.Request.RequestURI, c.Ctx.Input.RequestBody, resp)
}
