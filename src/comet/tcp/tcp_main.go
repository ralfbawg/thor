package tcp

import (
	"encoding/json"
	"net"
	"github.com/panjf2000/ants"
	"common/logging"
	"context"
	"bytes"
	"errors"
	"util"
)

const (
	TCP_MSG_TYPE_NORMAL    = iota
	TCP_MSG_TYPE_SUB
	TCP_MSG_TYPE_BROADCAST
	TCP_MSG_TYPE_UNICAST
	TCP_MSG_TYPE_PING      = 100
	TCP_MSG_TYPE_PONG      = 101
	TCP_MSG_TYPE_CLOSE     = 110

	BYTE_SIZE_COUNT  = 4
	BYTE_SIZE_SMALL  = 1024 * 1
	BYTE_SIZE_MEDIUM = 1024 * 32
	BYTE_SIZE_LARGE  = 1024 * 128
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
	TcpManagerInst.unbindClients.Set(ip, tcpClient)
	ants.Submit(tcpClient.run)
}
func wrapConc(conn net.Conn, ip string) *TcpClient {
	logging.Debug("tcp client connected,ip(%s)", ip)
	return &TcpClient{
		conn: conn,
		ip:   ip,
		send: make(chan []byte, 100),
		read: make(chan []byte, 100),
		c: TcpClientContext{
			context.TODO(),
		},
		//bytesPool: &sync.Pool{
		//	New: func() interface{} {
		//		return make([]byte, BYTE_SIZE_SMALL)
		//	},
		//},
	}
}

func MainEntrance(conn net.Conn, ip string) {
	HandlerConn(conn, ip)
}

func SendMsg(appId string, taskId int, uid string, msg []byte) error {
	a, exist := TcpManagerInst.bindClients.Get(appId)
	if !exist {
		return errors.New("there is no biz sub this app(" + appId + ")")
	}
	logging.Debug("appId %s bind tcp clients contains keys (%s)", appId, a.(util.ConcMap).Keys())
	if c, ok := TcpManagerInst.GetTcpClient(appId); ok {
		//if _, ok := bindClients.Get(appId); !ok {
		msg = bytes.Trim(msg, "\"")
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
			logging.Debug("send to tcp client(%s) msg(%s)", c.ip, resultB)
		} else {
			return errors.New("json解析出错，确定你的json没问题?")
		}

	}
	return nil
}
