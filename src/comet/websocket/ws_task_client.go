package websocket

import (
	"bytes"
	"common/logging"
	"github.com/gorilla/websocket"
	"time"
)

type WsTaskClient struct {
	task *WsTask

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	uid    string
	inited bool //注册，初始化完成
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3600 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 10

	hiMesaage     = "hi"
	hiMesaageJson = "{\"status\":\"ok\",\"uid\":\"%s\"}"
	helloMessage  = "hello"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *WsTaskClient) Close() {
	WsListenersInst.TriggerEvent(c.task.app.appId, c.uid, WS_EVENT_CLOSE)
	logging.Debug("close client uid(%s) by hand", c.uid)
	c.task.unregister <- c
	c.conn.Close()
}
func (c *WsTaskClient) readGoroutine() {
	WsListenersInst.TriggerEvent(c.task.app.appId, c.uid, WS_EVENT_CONNECTED)
	defer func() {
		WsListenersInst.TriggerEvent(c.task.app.appId, c.uid, WS_EVENT_CLOSE)
		logging.Debug("client uid(%s) closing", c.uid)
		c.task.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.conn.SetPingHandler(func(appData string) error {
		c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logging.Info("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if c.inited {

		}
		WsListenersInst.TriggerEvent(c.task.app.appId, c.uid, WS_EVENT_READ, message)
		logging.Debug("receive the websocket client(%s),ip(%s) message:%s ", c.uid, c.GetConn().RemoteAddr().String(), message)
		//c.send <- []byte(helloMessage)
		c.task.app.processMsg(0, c.uid, message)
		//logging.Debug("id %s get msg: %s",c.id,message)
	}
}

func (c *WsTaskClient) writeGoroutine() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The manager closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			WsListenersInst.TriggerEvent(c.task.app.appId, c.uid, WS_EVENT_WRITE, message)
			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WsTaskClient) Send(msg []byte) {
	if c != nil {
		c.send <- msg
	}

}
func (c *WsTaskClient) GetConn() *websocket.Conn {
	return c.conn
}
func isData(frameType int) bool {
	return frameType == websocket.TextMessage || frameType == websocket.BinaryMessage
}
