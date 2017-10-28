package models

type RobotReq struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserID string `json:"userid"`
}

type RobotResp struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func (r RobotResp) IsValid() bool {
	if r.Code == 100000 {
		return true
	}

	return false
}
