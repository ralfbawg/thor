package websocket

import (
	"github.com/panjf2000/ants"
	"strconv"
	"util"
)

const (
	MAP_KEY_SEPARATOR  = "#"
	WS_EVENT_CONNECTED = iota
	WS_EVENT_READ
	WS_EVENT_WRITE
	WS_EVENT_CLOSE
)

type WsListeners struct {
	util.ConcMap
}

type WsListenerI interface {
	TriggerEvent(event int)
	Register(f func(a ...interface{}))
}

type WsListener struct {
	WsListenerI
	appId string
	funcs [][]func(i ...interface{})
	eventChan chan WsEvent
}

func (l *WsListener) run() {
	select {
	case event := <-l.eventChan:
		registerFuncs := l.funcs[event.event]
		for _, eventFunc := range registerFuncs {
			ants.Submit(func() {
				eventFunc(event.param)
			})
		}
	}
}

type WsEvent struct {
	event int
	param []interface{}
}

func NewWsListener(appId string) *WsListener {
	return &WsListener{
		appId:     appId,
		eventChan: make(chan WsEvent),
	}
}
func (l WsListeners) TriggerEvent(appId string, event int, ext ...interface{}) {
	var funcs []func(i ...interface{})
	if tmp, ok := l.Get(appId); ok {
		funcs = tmp.([]func(i ...interface{}))
	}
	switch event {
	case WS_EVENT_CONNECTED, WS_EVENT_CLOSE:
		for _, fun := range funcs {
			fun()
		}
	case WS_EVENT_READ, WS_EVENT_WRITE:
		for _, fun := range funcs {
			fun(ext)
		}
	}
}
func (l *WsListeners) Register(appId string, event int, f ...func(a ...interface{})) {
	if tmp, ok := l.Get(appId + MAP_KEY_SEPARATOR + strconv.Itoa(event)); ok {
		funcs := tmp.([]func(i ...interface{}))
		f = append(funcs, f...) //TODO 同步问题
	}
	l.Set(appId, f)

}
