package tcp

import (
	"github.com/panjf2000/ants"
	"util"
	"bytes"
	"common/logging"
	"io"
	"encoding/json"
	"sync"
	"time"
)

const (
	WS_EVENT_CONNECTED = iota
	WS_EVENT_READ
	WS_EVENT_WRITE
	WS_EVENT_CLOSE
)

var (
	tcpCPoolExtendFactor = 0.8
	tcpCPoolDefaultSize  = 10000
	tcpCPool, _          = ants.NewPool(tcpCPoolDefaultSize)
	newline              = []byte{'\n'}
	space                = []byte{' '}
	singalByte           = make([]byte, 1)[0]
)

const (
	READ_TIME_OUT    = 60 * time.Second
	WRITE_TIME_OUT   = 60 * time.Second
	BYTE_SIZE_SMALL  = 1024 * 1
	BYTE_SIZE_MEDIUM = 1024 * 32
	BYTE_SIZE_LARGE  = 1024 * 128
)

var smallBytePool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, BYTE_SIZE_SMALL)
	},
}
var mediumBytePool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, BYTE_SIZE_MEDIUM)
	},
}
var largeBytePool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, BYTE_SIZE_LARGE)
	},
}

func (c *TcpClient) run() {
	util.SubmitTaskAndResize(tcpCPool, tcpCPoolDefaultSize, tcpCPoolExtendFactor, c.Write, c.Read)
	var resultB = []byte{}
	for {
		select {
		case msg := <-c.read:
			exist := bytes.Contains(msg, newline)
			logging.Debug("get tcp chan msg %s contain LR (%v)", msg, exist)
			if exist {
				lastIndex := bytes.LastIndex(msg, newline)
				resultB = append(resultB, msg[:lastIndex]...)
				for _, v := range bytes.Split(resultB, newline) {
					if len(v) > 0 {
						c.trimAndProcessMsg(v)
					}
				}
				if lastIndex+1 != len(msg) {
					errClone := util.Clone(msg[lastIndex+1:], &resultB)
					if errClone != nil {
						logging.Debug("%v clone error happen", errClone)
					}
				} else {
					resultB = []byte{}
				}

			} else {
				resultB = append(resultB, msg...)
			}
		}
	}

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
		b := make([]byte, BYTE_SIZE_SMALL)
		logging.Debug("buffer get by pool length is(%d),content is  (%s)", len(b), string(b))
		//b := smallBytePool.Get().([]byte)
		n, err := c.conn.Read(b)
		if err != nil {
			logging.Debug("tcp ip(%s) read error", c.ip)
			if err == io.EOF {
				logging.Error("got %v,tcp(%s) close ", err, c.ip)
				c.close()
				break
			} else {
				logging.Error("got %v; want %v", err, io.EOF)
				c.close()
				break
			}
		}
		c.read <- b[:n]
		logging.Debug("buffer read count is(%d) and length is(%d),content is  (%s)", n, len(b), string(b))
		buffer:=bytes.NewBuffer(b)
		buffer.Reset()
		logging.Debug("buffer after reset get by pool length is(%d),content is  (%s)", len(b), string(b))
		buffer.Bytes()
		smallBytePool.Put(b)

		//logging.Debug("(%s) contain LR(%v)", b, bytes.Contains(b, newline))
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
				//task.VerifyAppInfo2(reqMsg.Header.AppId, reqMsg.Header.TaskId, reqMsg.Header.Uid, reqMsg.Header.AppKey)
				c.appId = reqMsg.Header.AppId
				c.taskId = reqMsg.Header.TaskId
				c.uid = reqMsg.Header.Uid
				TcpManagerInst.WsListenerRegister(reqMsg.Header.AppId, c.ip, func(i ...interface{}) {
					//fmt.Printf("good event %s", i)
					//logging.Debug("i am trigger event")
					tcpMsg := &TcpMsg{}
					tcpMsg.Header.AppId = c.appId
					tcpMsg.Header.TaskId = c.taskId
					tcpMsg.Header.Uid = i[1].(string)
					tcpMsg.Header.MsgType = 200
					tcpMsg.Body = make(map[string]interface{}, 1)
					tcpMsg.Body["event"] = i[0]
					msg, _ := json.Marshal(tcpMsg)
					c.sendMsg(msg)
					//}, WS_EVENT_CONNECTED, WS_EVENT_READ, WS_EVENT_WRITE, WS_EVENT_CLOSE)
				}, WS_EVENT_CONNECTED, WS_EVENT_CLOSE)
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
		return nil, err
	}

	backMsg, err := json.Marshal(reqMsg)

	if err == nil {
		c.sendMsg(backMsg)
	} else {
		logging.Error("出错了")
	}
	return backMsg, err
}
func (c *TcpClient) close() {
	TcpManagerInst.WsListenerUnregister(c.appId, c.ip)
	TcpManagerInst.close <- c
	c.conn.Close()
}
func (c *TcpClient) sendMsg(msg []byte) {
	c.send <- msg
}
func (c *TcpClient) closeWs(uid string) {
	TcpManagerInst.WsCloseHandler(c.appId, 0, uid)
}
func (c *TcpClient) trimAndProcessMsg(msg []byte) error {
	msg = bytes.TrimSpace(msg)
	logging.Debug("get tcp message %s from %s", string(msg), c.ip)
	_, err := c.ProcessTcpMsg(msg)
	return err
}
