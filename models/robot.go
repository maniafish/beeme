package models

import (
	"beeme/util/counter"
	"errors"

	"github.com/astaxie/beego"
)

// RobotReq request to tuling api
type RobotReq struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserID string `json:"userid"`
}

// RobotResp response from tuling api
type RobotResp struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// IsValid check resp valid
func (r RobotResp) IsValid() bool {
	if r.Code == 100000 {
		return true
	}

	return false
}

// IsLimit reach max request limit in one day
func (r RobotResp) IsLimit() bool {
	if r.Code == 40004 {
		return true
	}

	return false
}

var (
	robotid = counter.New()
	// ErrRML reach max limit error
	ErrRML = errors.New("Reach max Limit")
)

// GetRobotKey choose one of Tuling robot
func GetRobotKey() string {
	keyID := robotid.Incr() % uint64(len(beego.AppConfig.Strings("apps::TulingKeys")))
	return beego.AppConfig.Strings("apps::TulingKeys")[keyID]
}

const (
	// DefaultReturn default return to wechat when internal err
	DefaultReturn = "success"
	// DefaultErrMsg msg when program err
	DefaultErrMsg = "小O有点转不过弯来了~"
	// SubscribeMsg msg of subscribe return
	SubscribeMsg = "Once in a Life. 一生一次，一次一生～小O(Oil)愿为您服务至过劳～回复“呼叫小O”，告诉您小O的故事，以及一些好玩的功能~"
	// RmlErrMsg msg when reach request limit in one day
	RmlErrMsg = "小O今天要休息啦，客官明天再来吧～"
)

// RobotMsg robot msg model
type RobotMsg struct {
	ID      int    `orm:"pk;column(id);auto"`
	Msg     string `orm:"unique"`
	MsgType string `orm:"column(msg_type)"`
	Resp    string
}

// Get get resp by msg
func (r *RobotMsg) Get() error {
	return userOrmer.Raw("SELECT id, msg, msg_type, resp FROM robot_msg WHERE msg = ?", r.Msg).QueryRow(r)
}

// TableIndex set index
func (r *RobotMsg) TableIndex() [][]string {
	return [][]string{
		[]string{"MsgType", "Resp"},
	}
}
