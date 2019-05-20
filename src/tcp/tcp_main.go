package tcp

import (
	"encoding/json"
	"net"
	"github.com/panjf2000/ants"
)

const (
	CONNECT_TYPE_SHORT = iota
	CONNECT_TYPE_LONG  
)

func ProcessTcpMsg(msg []byte) *TcpMsg {
	reqMsg := &TcpMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		switch reqMsg.Header.MsgType {
		//0:订阅 1:广播 2:单播 110:断掉客户端长连接
		case 0:
		case 1:
		case 2:
		case 110:

		}
	}
	return nil
}
func UnmarshalMsg(msg []byte) (*TcpMsg, error) {
	reqMsg := &TcpMsg{}
	json.Marshal(reqMsg.Header)
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		return reqMsg, nil
	} else {
		return nil, err
	}

}
func HandlerConn(conn net.Conn, ip string) {
	tcpClient := wrapConc(conn, ip)
	unbindClients.Set(ip, tcpClient)
	ants.Submit(tcpClient.run)
}
func wrapConc(conn net.Conn, ip string) *TcpClient {
	return &TcpClient{
		conn: conn,
		ip:   ip,
	}
}

func MainEntrance(conn net.Conn, ip string) {
	HandlerConn(conn, ip)
}

func good() {

}
