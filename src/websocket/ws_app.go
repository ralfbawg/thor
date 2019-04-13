package websocket

import "util"

type WsApp struct {
	AppId string
	Tasks util.ConcMap
}
