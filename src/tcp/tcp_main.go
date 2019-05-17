package tcp

import (
	"encoding/json"
	"net"
	"io"
	"common/logging"
	"bytes"
	"util"
	"github.com/panjf2000/ants"
)

const (
	CONNECT_TYPE_SHORT = iota
	CONNECT_TYPE_LONG
)

var (
	tcpCPoolExtendFactor       = 0.8
	tcpCPoolDefaultSize        = 10000
	tcpCPool, _                = ants.NewPool(tcpCPoolDefaultSize)
	funcs                      = make([]func(), 0)
	newline                    = []byte{'\n'}
	space                      = []byte{' '}
	bindClients, unbindClients = util.NewConcMap(), util.NewConcMap()
)

func (c *TcpClient) run() {
	util.SubmitTaskAndResize(tcpCPool, tcpCPoolDefaultSize, tcpCPoolExtendFactor, append(funcs[:0], c.Write, c.Read))
}

func (c *TcpClient) Write(msg []byte) {
	for {
		select {
			msg := <-c.send
			c.conn.
		}
	}
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

func ProcessTcpMsg(msg []byte) *TcpMsg {
	reqMsg := &TcpMsg{}
	err := json.Unmarshal(msg, reqMsg)
	if err == nil {
		return reqMsg
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
func HandlerConn(conn *net.Conn, ip string) {
	tcpClient := wrapConc(conn, ip)
	unbindClients.Set(ip, tcpClient)
}
func wrapConc(conn *net.Conn, ip string) *TcpClient {
	return &TcpClient{
		conn: conn,
		ip:   ip,
	}
}

func MainEntrance(conn net.Conn, ip string) {

}
