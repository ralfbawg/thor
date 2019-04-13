package game

import (
	"bytes"
	"common/logging"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants"
	"time"
)

const (
	ROOM_POS_EMPTY = -1
	ROOM_POS_A     = iota
	ROOM_POS_B
	ROOM_POS_ALL     //both
	USER_EVENT_START = "start"
	USER_EVENT_EXIT  = "exit"
)

type GameClient struct {
	gm *GameMall

	gameRoom *GameRoom

	opp *GameClient
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	read chan []byte

	id string

	pos int32
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *GameClient) run() {
	ants.Submit(c.readGoroutine)
	ants.Submit(c.writeGoroutine)
}

func (c *GameClient) findGame() {
	c.gm.findClientId <- c.id
}
func (c *GameClient) exitGame() {
	c.gameRoom.ExitClient(c, false)
}
func (c *GameClient) closeGame() {
	c.gameRoom.ExitClient(c, true)
}
func (c *GameClient) readGoroutine() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.conn.SetPingHandler(func(appData string) error {
		c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(pongWait))
		return nil
	})
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.closeGame()
		if c != nil && c.conn != nil {
			message := websocket.FormatCloseMessage(code, "")
			c.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(writeWait))
		}
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
		logging.Info("get message %s from client(%s)", message, c.id)
		if string(message) == USER_EVENT_START {
			c.findGame()
		} else if string(message) == USER_EVENT_EXIT {
			c.exitGame()
		} else {
			c.read <- message
		}
	}
}

func (c *GameClient) writeGoroutine() {
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
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

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

func (c *GameClient) Send(msg []byte) {
	if c != nil {
		c.send <- msg
	}

}

func (c *GameClient) ID() string {
	return c.id
}
func (c *GameClient) IP() string {
	return c.conn.RemoteAddr().String()
}
