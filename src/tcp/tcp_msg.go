package tcp

import "net"

type TcpMsg struct {
	Header TcpMsgHeader `json:"header"`
	Body   TcpMsgBody   `json:"body"`
}
type TcpMsgHeader struct {
	MsgType int    `json:"type"`
	AppId   string `json:"appId"`
	AppKey  string `json:"appKey"`
	TaskId  string `json:"taskId"`
	Uid     string `json:"uid"`
}
type TcpMsgBody struct {
	Uri  string `json:"uri"`
	Data []byte `json:"data"`
}
type TcpClient struct {
	ConnectType string
	conn        net.Conn
	appId       string
	taskId      int64
	uid         string
	ip          string
	send        chan []byte
	read        chan []byte
}
