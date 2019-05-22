package websocket

import "util"

const (
	WS_EVENT_CONNECTED = iota
	WS_EVENT_READ
	WS_EVENT_WRITE
	WS_EVENT_CLOSE
)

type WsListenerI interface {
	OnEvent(event int)
	Register(f func(a interface{}))
}

type BaseWsListener struct {
	WsListenerI
	listeners util.ConcMap
}

func NewWsListener() *BaseWsListener {
	return &BaseWsListener{
		listeners: util.NewConcMap(),
	}
}
func (l *BaseWsListener) OnEvent(event int) {
	switch event {
	case WS_EVENT_CONNECTED:
	case WS_EVENT_READ:
	case WS_EVENT_WRITE:
	case WS_EVENT_CLOSE:

	}
}
func (l *BaseWsListener) Register(appId string, f func(a interface{})) {

	if tmp, ok := l.listeners.Get(appId);ok {
		tmpA := tmp.([]func)
	}
}
