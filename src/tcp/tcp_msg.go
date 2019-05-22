package tcp

import "net"

type TcpMsg struct {
	Header *TcpMsgHeader          `json:"header"`
	Body   map[string]interface{} `json:"body"`
}
type TcpMsgHeader struct {
	MsgType   int    `json:"type"`
	AppId     string `json:"appId"`
	AppKey    string `json:"appKey"`
	TaskId    int    `json:"taskId"`
	Uid       string `json:"uid"`
	RequestId string `json:"requestId"`
}

//type TcpMsgBody struct {
//	Uri  string `json:"uri"`
//	Data string `json:"data"`
//}
type TcpClient struct {
	ConnectType string
	conn        net.Conn
	appId       string
	taskId      int
	uid         string
	ip          string
	send        chan []byte
	read        chan []byte
}
