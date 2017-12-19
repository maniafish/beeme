package hello

import (
	"beeme/util/xmls"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
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
		c.ServeString("success")
		return
	}

	resp := &MessageResp{
		ToUserName:   r.ToUserName,
		FromUserName: r.FromUserName,
		MsgType:      r.MsgType,
		CreateTime:   time.Now().Unix(),
	}

	switch msgType := strings.ToLower(string(r.MsgType)); msgType {
	case "text":
		resp.Content = xmls.CharData("小O来咯~")
		c.ServeXML(resp)
		return
	case "image", "voice", "video", "shortvideo", "location", "link":
		resp.Content = xmls.CharData(fmt.Sprintf("小O还搞不定%s哦~", msgTypeMap[msgType]))
		c.ServeXML(resp)
		return
	default:
		c.Logger.Infof("invalid MsgType: %v", msgType)
		c.ServeString("success")
		return
	}
}
