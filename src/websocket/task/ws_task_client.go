package task

import "github.com/gorilla/websocket"

type WsTaskClient struct {
	hub *WsTask

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	id string
}