package tcp

import (
	"encoding/json"
	"net"
	"io"
	"common/logging"
	"bytes"
	"util"
)

const (
	CONNECT_TYPE_SHORT = iota
	CONNECT_TYPE_LONG
)

var (
	newline       = []byte{'\n'}
	space         = []byte{' '}
	bindClients   = util.NewConcMap()
	unbindClients = util.NewConcMap()
)

type TcpMsg struct {
	appId string
	uri   string
	uid   string
	data  []byte
}
type TcpClient struct {
	ConnectType string
	conn        net.Conn
	appId       string
	taskId      int64
	uid         string
}

func (c *TcpClient) Write(msg []byte) {
	c.conn.Write(msg)
}
func (c *TcpClient) Read() []byte {
	var message [1]byte
	for {
		if _, err := c.conn.Read(message[:]); err != io.EOF {
			logging.Error("got %v; want %v", err, io.EOF)
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//logging.Debug("the msg type is %d", msgType)
		c.send <- []byte(helloMessage)

		//logging.Debug("id %s get msg: %s",c.id,message)
	}
	c.conn.Write(msg)

	c.conn.Read(msgB)
	return msgB
}

func ProcessTcpMsg(msg []byte) {
	reqMsg := &TcpMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		return
	}

}
func UnmarshalMsg(msg []byte) (*TcpMsg, error) {
	reqMsg := &TcpMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		return reqMsg, nil
	} else {
		return nil, err
	}

}
func HanlderConc(conn net.Conn) {
	client := wrapConc(conn)
}
func wrapConc(conn net.Conn) *TcpClient {
	return &TcpClient{
		conn: conn,
	}
}
