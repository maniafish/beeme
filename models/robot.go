package models

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
