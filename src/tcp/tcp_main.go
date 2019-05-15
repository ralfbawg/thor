package tcp

import "encoding/json"

type ReqMsg struct {
}

func ProcessTcpMsg(msg []byte) {
	reqMsg := &ReqMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		return
	}

}
func UnmarshalMsg(msg []byte) (*ReqMsg, error) {
	reqMsg := &ReqMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		return reqMsg, nil
	} else {
		return nil, err
	}

}
