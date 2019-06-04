package websocket

import "util"

const (
	WS_EVENT_CONNECTED = iota
	WS_EVENT_READ
	WS_EVENT_WRITE
	WS_EVENT_CLOSE
)

var Wslisteners = NewWsListener()

type WsListenerI interface {
	OnEvent(event int)
	Register(f func(a ...interface{}))
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
func (l *BaseWsListener) TriggerEvent(appId string, event int, ext ...interface{}) {
	funcs := make([]func(i ...interface{}), 0)
	if tmp, ok := l.listeners.Get(appId); ok {
		funcs = tmp.([]func(i ...interface{}))
	}
	switch event {
	case WS_EVENT_CONNECTED, WS_EVENT_CLOSE:
		for _, fun := range funcs {
			fun(event)
		}
	case WS_EVENT_READ, WS_EVENT_WRITE:
		for _, fun := range funcs {
			fun(event, ext)
		}
	}
}
func (l *BaseWsListener) Register(appId string, f ...func(a ...interface{})) {
	if tmp, ok := l.listeners.Get(appId); ok {
		funcs := tmp.([]func(i ...interface{}))
		f = append(funcs, f...)
	}
	l.listeners.Set(appId, f)

}
