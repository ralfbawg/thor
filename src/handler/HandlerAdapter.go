package handler

import (
	"github.com/gorilla"
	"time"
)

const (
	heartBeat         = "hb"
	heartBeatInterval = 1 * time.Second
)

func AdaptTask(conn websocket.Conn) (string, error) {
	var (
		err error
	)
	for {
		if err = conn.WriteMessage([]byte(heartBeat)); err != nil {
			return
		}
		time.Sleep(heartBeatInterval)
	}
	go adaptTask(&conn)
}

func adaptTask(conn *websocket.Conn)  (string, error)  {
	var (
		err error
	)
	for {
		if data, err = *conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = *conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}

