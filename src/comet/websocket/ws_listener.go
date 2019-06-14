package websocket

import (
	"common/logging"
	"context"
	"github.com/panjf2000/ants"
	"util"
	"sync"
)

const (
	WS_EVENT_CONNECTED = iota
	WS_EVENT_READ
	WS_EVENT_WRITE
	WS_EVENT_CLOSE
	MAP_KEY_SEPARATOR  = "#"
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
	for {
		select {
		case event := <-l.eventChan:
			registerFuncs := l.funcs[event.event]
			for _, eventFunc := range registerFuncs {
				ants.Submit(func() {
					eventFunc(event.event, event.param)
				})
			}
		case <-l.c.Done():
			logging.Debug("WsListener parent context done,i should exit")
			close(l.eventChan) //不再接收事件
			return
		}
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
		eventChan: make(chan *WsEvent, 10),
		c:         context,
		funcs: make([][]func(params ...interface{}), 10),
	}
}
func (l *WsListeners) TriggerEvent(appId string, event int, ext ...interface{}) {
	var listener *WsListener
	if tmp, ok := l.Get(appId); ok {
		listener = tmp.(*WsListener)
		listener.eventChan <- &WsEvent{event: event, param: ext,}
	} else {
		logging.Debug("appId %s is not exist", appId)
	}

}

//注册监听事件
func (l *WsListeners) Register(appId string, f func(a ...interface{}), events ...int) {
	l.lock.RLock() //TODO 同步问题有没有更好的方法呢?
	defer l.lock.RUnlock()
	var listener *WsListener
	if tmp, ok := l.Get(appId); ok {
		listener = tmp.(*WsListener)

	} else {
		listener = NewWsListener(WsListenersInst.context, appId)
		ants.Submit(listener.run)
	}
	for _, event := range events {
		if listener.funcs[event] == nil {
			listener.funcs[event] = []func(params ...interface{}){
				f,
			}
		} else {
			listener.funcs[event] = append(listener.funcs[event], f)
		}
	}
	l.Set(appId, listener)
}
