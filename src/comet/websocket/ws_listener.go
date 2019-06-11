package websocket

import (
	"common/logging"
	"context"
	"github.com/panjf2000/ants"
	"util"
	"sync"
)

const (
	MAP_KEY_SEPARATOR  = "#"
	WS_EVENT_CONNECTED = iota
	WS_EVENT_READ
	WS_EVENT_WRITE
	WS_EVENT_CLOSE
)

var WsListenersInst = &WsListeners{
	ConcMap: util.NewConcMap(),
	context: context.Background(),
}

type WsListeners struct {
	util.ConcMap
	context context.Context
	lock    sync.RWMutex
}

type WsListenerI interface {
	TriggerEvent(event int)
	Register(f func(a ...interface{}))
}

type WsListener struct {
	WsListenerI
	appId string
	funcs [][]func(i ...interface{})
	eventChan chan *WsEvent
	c         context.Context
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
	case <-l.c.Done():
		logging.Debug("WsListener parent context done,i should exit")
		close(l.eventChan) //不再接收事件
		return
	}
}

type WsEvent struct {
	event int
	param []interface{}
}

func NewWsListener(c context.Context, appId string) *WsListener {
	context, _ := context.WithCancel(c)
	return &WsListener{
		appId:     appId,
		eventChan: make(chan *WsEvent),
		c:         context,
		funcs: make([][]func(params ...interface{}), 10),
	}
}
func (l WsListeners) TriggerEvent(appId string, event int, ext ...interface{}) {
	var listner *WsListener
	if tmp, ok := l.Get(appId); ok {
		listner = tmp.(*WsListener)
		listner.eventChan <- &WsEvent{event: event, param: ext,}
	} else {
		logging.Debug("appId %s is not exist", appId)
	}

}

//注册监听事件
func (l *WsListeners) Register(appId string, event int, f ...func(a ...interface{})) {
	l.lock.RLock() //TODO 同步问题有没有更好的方法呢?
	defer l.lock.RUnlock()
	if tmp, ok := l.Get(appId); ok {
		ll := tmp.(*WsListener)

		if ll.funcs[event] == nil {
			ll.funcs[event] = make([]func(params ...interface{}), 10)
		} else {
			ll.funcs[event] = append(ll.funcs[event], f...)
		}
	} else {
		listner := NewWsListener(WsListenersInst.context, appId)
		ants.Submit(listner.run)
		l.Set(appId, listner)
	}

}
