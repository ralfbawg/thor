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
	Context: context.Background(),
}

type WsListeners struct {
	util.ConcMap
	Context context.Context
	lock    sync.RWMutex
}

type WsListenerI interface {
	TriggerEvent(event int)
	Register(f func(a ...interface{}))
}

type WsListener struct {
	WsListenerI
	appId string
	//eventFunctions [][]func(i ...interface{})
	eventFunctions util.ConcMap
	eventChan      chan *WsEvent
	c              context.Context
}

func (l *WsListener) run() {
	for {
		select {
		case event := <-l.eventChan:
			for tup := range l.eventFunctions.IterBuffered() {
				registerFuncs := tup.Val.([][]func(i ...interface{}))[event.event]
				for _, eventFunc := range registerFuncs {
					ants.Submit(func() {
						eventFunc(event.event, event.uid, event.param)
					})
				}
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
	uid   string
	param []interface{}
}

func NewWsListener(c context.Context, appId string) *WsListener {
	context, _ := context.WithCancel(c)
	return &WsListener{
		appId:          appId,
		eventChan:      make(chan *WsEvent, 10),
		c:              context,
		eventFunctions: util.NewConcMap(),
	}
}
func (l *WsListeners) TriggerEvent(appId string, uid string, event int, ext ...interface{}) {
	if tmp, ok := l.Get(appId); ok {
		event := &WsEvent{event: event, uid: uid, param: ext,}
		listener := tmp.(*WsListener)
		listener.eventChan <- event
	} else {
		logging.Debug("appId %s is not exist", appId)
	}

}

//注册监听事件
func (l *WsListeners) Register(appId string, ip string, f func(a ...interface{}), events ...int) {
	l.lock.RLock() //TODO 同步问题有没有更好的方法呢?
	defer l.lock.RUnlock()
	var listener *WsListener
	var c []func(i ...interface{})
	if tmp, ok := l.Get(appId); ok {
		listener = tmp.(*WsListener)
	} else {
		listener = NewWsListener(WsListenersInst.Context, appId)
		ants.Submit(listener.run)
	}
	if a, b := listener.eventFunctions.Get(ip); b {
		c = a.([]func(i ...interface{}))
	} else {
		c = make([]func(i ...interface{}), 10)
	}
	for _, event := range events {
		c[event] = f
	}
	l.Set(appId, listener)
}

//注册监听事件
func (l *WsListeners) Unregister(appId string, ip string) {
	if a, ok := l.Get(appId); ok {
		if b, okB := a.(util.ConcMap).Get(ip); okB {
			close(b.(*WsListener).eventChan)
		}
		a.(util.ConcMap).Remove(ip)
	}

}
