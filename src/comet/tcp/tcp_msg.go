package tcp

import (
	"net"
	"context"
)

type TcpMsg struct {
	Header struct {
		MsgType   int    `json:"type"`
		AppId     string `json:"appId"`
		AppKey    string `json:"appKey"`
		TaskId    int    `json:"taskId"`
		Uid       string `json:"uid"`
		RequestId string `json:"requestId"`
	} `json:"header"`
	Body map[string]interface{} `json:"body"`
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
	c           TcpClientContext
}

type TcpClientContext struct {
	context.Context
}
