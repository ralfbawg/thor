package tcp

import (
	"encoding/json"
	"net"
	"github.com/panjf2000/ants"
	"common/logging"
)

const (
	TCP_MSG_TYPE_NORMAL    = iota
	TCP_MSG_TYPE_SUB
	TCP_MSG_TYPE_BROADCAST
	TCP_MSG_TYPE_UNICAST
	TCP_MSG_TYPE_PING      = 100
	TCP_MSG_TYPE_PONG      = 101
	TCP_MSG_TYPE_CLOSE     = 110
)

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
		send: make(chan []byte, 100),
		read: make(chan []byte, 100),
	}
}

func MainEntrance(conn net.Conn, ip string) {
	HandlerConn(conn, ip)
}

func SendMsg(appId string, taskId int, uid string, msg []byte) {
	a := bindClients.Keys()
	logging.Info("%s", a)
	if v, ok := bindClients.Get(appId); ok {
		//if _, ok := bindClients.Get(appId); !ok {
		c := v.(*TcpClient)

		tcpMsg := &TcpMsg{
			Body: make(map[string]interface{}),
		}
		json.Unmarshal(msg, &tcpMsg.Body)
		tcpMsg.Header.AppId = appId
		tcpMsg.Header.TaskId = taskId
		tcpMsg.Header.Uid = uid
		resultB, err := json.Marshal(tcpMsg)
		if err == nil {
			c.send <- resultB
		} else {
			c.send <- []byte("哦活，出错了")
		}
		logging.Info("%s", resultB)
	}
}
