package wechat

import (
	"beeme/models"
	"beeme/util/xmls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
)

// MessageReq request for wechat message
type MessageReq struct {
	XMLName      xml.Name      `xml:"xml" sort:"-"`
	ToUserName   xmls.CharData `xml:"ToUserName" sort:"ToUserName,omitempty"`
	FromUserName xmls.CharData `xml:"FromUserName" sort:"FromUserName,omitempty"`
	CreateTime   int64         `xml:"CreateTime" sort:"CreateTime,omitempty"`
	MsgType      xmls.CharData `xml:"MsgType" sort:"MsgType"`
	MsgID        int64         `xml:"MsgId" sort:"MsgId,omitempty"`
	Content      xmls.CharData `xml:"Content" sort:"Content,omitempty"`           // for MsgType: text
	Event        xmls.CharData `xml:"Event" sort:"Event,omitempty"`               // for MsgType: event
	MediaID      xmls.CharData `xml:"MediaId" sort:"MediaId,omitempty"`           // for MsgType: image/voice/video/shortvideo, pull data from server
	PicURL       xmls.CharData `xml:"PicUrl" sort:"PicUrl,omitempty"`             // for MsgType: image
	Format       xmls.CharData `xml:"Format" sort:"Format,omitempty"`             // for MsgType: voice
	ThumbMediaID xmls.CharData `xml:"ThumbMediaId" sort:"ThumbMediaId,omitempty"` // for MsgType: video/shortvideo
	LocationX    xmls.CharData `xml:"Location_X" sort:"Location_X,omitempty"`     // for MsgType: location
	LocationY    xmls.CharData `xml:"Location_Y" sort:"Location_Y,omitempty"`     // for MsgType: location
	Scale        xmls.CharData `xml:"Scale" sort:"Scale,omitempty"`               // for MsgType: location
	Label        xmls.CharData `xml:"Label" sort:"Scale,omitempty"`               // for MsgType: location
	Title        xmls.CharData `xml:"Title" sort:"Title,omitempty"`               // for MsgType: link
	Description  xmls.CharData `xml:"Description" sort:"Description,omitempty"`   // for MsgType: link
	URL          xmls.CharData `xml:"Url" sort:"Url,omitempty"`                   // for MsgType: link
}

// MessageResp response for wechat
type MessageResp struct {
	XMLName      xml.Name      `xml:"xml"`
	ToUserName   xmls.CharData `xml:"ToUserName,omitempty"`
	FromUserName xmls.CharData `xml:"FromUserName,omitempty"`
	CreateTime   int64         `xml:"CreateTime"`
	MsgType      xmls.CharData `xml:"MsgType"`
	Content      xmls.CharData `xml:"Content"`
}

var msgTypeMap = map[string]string{
	"text":       "文本消息",
	"event":      "事件消息",
	"image":      "图片消息",
	"voice":      "语音消息",
	"video":      "视频消息",
	"shortvideo": "小视频消息",
	"location":   "地理位置消息",
	"link":       "链接消息",
}

// Post reutrn for wechat msg
// @Title Post
// @Description reutrn for wechat msg, doc site: https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140453
// @router / [post]
func (c *Controller) Post() {
	c.Initialize("wechat_message.Post()")
	r := &MessageReq{}
	err := xml.Unmarshal(c.Ctx.Input.RequestBody, r)
	if err != nil {
		c.Logger.Errorf("Unmarshal RequestBody err: %v", err)
		c.ServeString(models.DefaultReturn)
		return
	}

	resp := c.setDefaultReturn(r)
	switch msgType := strings.ToLower(string(r.MsgType)); msgType {
	case "text":
		err = c.textResponse(resp)
	case "event":
		err = c.eventResponse(resp)
	case "image", "voice", "video", "shortvideo", "location", "link":
		err = c.otherResponse(resp)
	default:
		err = errors.Errorf("Invalid MsgType: %v", msgType)
		c.Logger.Errorf(err.Error())
	}

	switch err {
	case nil:
	case models.ErrRML:
		resp.MsgType = xmls.CharData("text")
		resp.Content = xmls.CharData(models.RmlErrMsg)
	default:
		resp.MsgType = xmls.CharData("text")
		resp.Content = xmls.CharData(models.DefaultErrMsg)
	}

	c.ServeXML(resp)
	return
}

func (c *Controller) setDefaultReturn(r *MessageReq) *MessageResp {
	c.SetRequest(r)
	return &MessageResp{
		ToUserName:   r.FromUserName,
		FromUserName: r.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      xmls.CharData("text"),
	}
}

func (c *Controller) textResponse(resp *MessageResp) error {
	r, ok := c.GetRequest().(*MessageReq)
	if !ok {
		err := errors.New("GetRequest err")
		c.Logger.Errorf(err.Error())
		return err
	}

	// search db
	robotMsg := &models.RobotMsg{Msg: string(r.Content)}
	err := robotMsg.Get()
	switch err {
	case nil:
		if robotMsg.Resp != "" {
			resp.Content = xmls.CharData(robotMsg.Resp)
			return nil
		}
	case orm.ErrNoRows:
	default:
		c.Logger.Errorf("Get RobotMsg DBErr: %v", err)
		return err
	}

	// post Tuling
	req, _ := json.Marshal(&models.RobotReq{
		Key:    models.GetRobotKey(),
		Info:   string(r.Content),
		UserID: string(r.FromUserName),
	})

	respTuling, err := httplib.Post(beego.AppConfig.String("apps::TulingURL")).Body(req).Bytes()
	if err != nil {
		c.Logger.Errorf("Post Tulling err: %v", err)
		return err
	}

	obj := &models.RobotResp{}
	err = json.Unmarshal(respTuling, obj)
	if err != nil {
		c.Logger.Errorf("Unmarshal Resp err: %v", err)
		return err
	}

	if !obj.IsValid() {
		if obj.IsLimit() {
			return models.ErrRML
		}

		err = errors.Errorf("Invalid Resp: %+v", obj)
		c.Logger.Errorf(err.Error())
		return err
	}

	resp.Content = xmls.CharData(obj.Text)
	return nil
}

func (c *Controller) eventResponse(resp *MessageResp) error {
	r, ok := c.GetRequest().(*MessageReq)
	if !ok {
		err := errors.New("GetRequest err")
		c.Logger.Errorf(err.Error())
		return err
	}

	switch event := strings.ToLower(string(r.Event)); event {
	case "subscribe":
		resp.Content = xmls.CharData(models.SubscribeMsg)
		return nil
	default:
		err := errors.Errorf("Invalid Event: %v", event)
		c.Logger.Errorf(err.Error())
		return err
	}
}

func (c *Controller) otherResponse(resp *MessageResp) error {
	r, ok := c.GetRequest().(*MessageReq)
	if !ok {
		err := errors.New("GetRequest err")
		c.Logger.Errorf(err.Error())
		return err
	}

	resp.Content = xmls.CharData(fmt.Sprintf("小O还搞不定%s哦~", msgTypeMap[string(r.MsgType)]))
	return nil
}
