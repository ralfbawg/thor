package tcp

import (
	"github.com/panjf2000/ants"
	"util"
	"bytes"
	"common/logging"
	"io"
	"encoding/json"
	"task"
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

func (c *TcpClient) Write() {
	for {
		select {
		case msg := <-c.send:
			c.conn.Write(append(msg, newline...))
		}
	}
}
func (c *TcpClient) Read() {
	for {
		b := make([]byte, 256)
		n, err := c.conn.Read(b)
		if n != 0 {
			b = bytes.TrimSpace(b)
			bs := make([][]byte, 1)
			if bytes.Contains(b, newline) {
				bs = bytes.Split(b, newline)
			} else {
				bs[0] = b
			}
			for _, v := range bs {
				b = v[0:n]
				logging.Info("get tcp message %s from %s", string(b), c.ip)
				c.ProcessTcpMsg(b)
			}
			b = b[:0]
		} else if err == io.EOF {
			logging.Info("got %v; want %v", err, io.EOF)
			c.close()
			break
		} else {
			logging.Info("got %v; want %v", err, io.EOF)
			c.close()
			break
		}

	}

}
func (c *TcpClient) ProcessTcpMsg(msg []byte) ([]byte, error) {
	reqMsg := &TcpMsg{}
	err := json.Unmarshal(msg, reqMsg)

	if err == nil {
		switch reqMsg.Header.MsgType {
		//0:订阅 1:广播 2:单播 110:断掉客户端长连接 100心跳 101回应
		case TCP_MSG_TYPE_SUB:
			if reqMsg.Header.AppId != "" && reqMsg.Header.AppKey != "" {
				task.VerifyAppInfo2(reqMsg.Header.AppId, reqMsg.Header.TaskId, reqMsg.Header.Uid, reqMsg.Header.AppKey)
				c.appId = reqMsg.Header.AppId
				c.taskId = reqMsg.Header.TaskId
				c.uid = reqMsg.Header.Uid
				TcpManagerInst.bind <- c
				reqMsg.Header.MsgType = TCP_MSG_TYPE_PONG
			}
		case TCP_MSG_TYPE_BROADCAST:
		case TCP_MSG_TYPE_UNICAST:
			msg, _ := json.Marshal(reqMsg.Body)
			if c.appId != "" {
				TcpManagerInst.WsBroadcast(c.appId, c.taskId, c.uid, msg)
			}
		case TCP_MSG_TYPE_PING:
		case TCP_MSG_TYPE_PONG:
		case TCP_MSG_TYPE_CLOSE:

		}
	} else {
		logging.Error("[GetResponseInfo] Failed, %s,", err.Error())
	}

	backMsg, err := json.Marshal(reqMsg)

	if err == nil {
		c.send <- backMsg
	} else {
		logging.Error("出错了")
	}
	return backMsg, err
}
func (c *TcpClient) close() {
	c.conn.Close()
}
