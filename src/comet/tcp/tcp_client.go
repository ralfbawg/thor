package tcp

import (
	"comet/websocket"
	"github.com/panjf2000/ants"
	"util"
	"bytes"
	"common/logging"
	"io"
	"encoding/json"
	"task"
	"sync"
	"time"
)

var (
	tcpCPoolExtendFactor       = 0.8
	tcpCPoolDefaultSize        = 10000
	tcpCPool, _                = ants.NewPool(tcpCPoolDefaultSize)
	funcs                      = make([]func(), 0)
	newline                    = []byte{'\n'}
	space                      = []byte{' '}
	bindClients, unbindClients = util.NewConcMap(), util.NewConcMap()
	singalByte                 = make([]byte, 1)[0]
)

const (
	READ_TIME_OUT  = 60 * time.Second
	WRITE_TIME_OUT = 60 * time.Second
)

var bytePool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024*8)
	},
}

func (c *TcpClient) run() {
	util.SubmitTaskAndResize(tcpCPool, tcpCPoolDefaultSize, tcpCPoolExtendFactor, append(funcs[:0], c.Write, c.Read))
}

func (c *TcpClient) Write() {
	for {
		select {
		case msg := <-c.send:
			logging.Debug("write to tcp client appId=%s,content=%s", c.appId, msg)
			c.conn.Write(append(msg, newline...))
			//case c.c.Context.Done():
			//	return

		}
	}
}
func (c *TcpClient) Read() {
	for {
		b := bytePool.Get().([]byte)
		//b := make([]byte, 256)
		n, err := c.conn.Read(b)
		if n != 0 {
			b = bytes.TrimSpace(b)
			//var bs [][]byte
			////var bs  = [][]byte{b}
			//if bytes.Contains(b, newline) {
			//	bs = bytes.Split(b, newline)
			//} else {
			//	bs = [][]byte{b}
			//}
			//
			//for _, v := range bs {
			//	validN := bytes.IndexByte(v, singalByte)
			//	b = func() []byte {
			//		if validN < 0 {
			//			return v[0:]
			//		} else {
			//			return v[0:validN]
			//		}
			//	}()
			//
			//	logging.Info("get tcp message %s from %s", string(b), c.ip)
			//	c.ProcessTcpMsg(b)
			//}
			logging.Debug("get tcp message %s from %s", string(b[0:n]), c.ip)
			c.ProcessTcpMsg(b[0:n])
			bytes.NewBuffer(b).Reset()
			bytePool.Put(b)
		} else if err == io.EOF {
			logging.Error("got %v; want %v", err, io.EOF)
			c.close()
			break
		} else {
			logging.Error("got %v; want %v", err, io.EOF)
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
				websocket.Wslisteners.Register(reqMsg.Header.AppId, func(a ...interface{}) {
					
				})
				TcpManagerInst.bind <- c
			}
		case TCP_MSG_TYPE_BROADCAST, TCP_MSG_TYPE_UNICAST:
			msg, _ := json.Marshal(reqMsg.Body)
			if c.appId != "" {
				TcpManagerInst.WsBroadcast(c.appId, reqMsg.Header.TaskId, reqMsg.Header.Uid, msg)
			}
		case TCP_MSG_TYPE_PING, TCP_MSG_TYPE_PONG:
			logging.Debug("ping/pong")
			reqMsg.Header.MsgType = TCP_MSG_TYPE_PONG
		case TCP_MSG_TYPE_CLOSE:
			c.closeWs(reqMsg.Header.Uid)
		}
	} else {
		logging.Error("[ProcessTcpMsg] Failed, %s,", err.Error())
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
func (c *TcpClient) closeWs(uid string) {
	TcpManagerInst.WsCloseHandler(c.appId, 0, uid)
}
