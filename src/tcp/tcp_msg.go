package tcp

import "net"

type TcpMsg struct {
	heder TcpMsgHeader
	uri   string
	data  []byte
}
type TcpMsgHeader struct {
	appId  string
	taskId string
	uid    string
}

type TcpClient struct {
	ConnectType string
	conn        *net.Conn
	appId       string
	taskId      int64
	uid         string
	ip          string
}
