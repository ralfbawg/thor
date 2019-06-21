package tcp

import (
	"github.com/panjf2000/ants"
	"util"
	"time"
)

type TcpManager struct {
	bind           chan *TcpClient
	unbind         chan *TcpClient
	close          chan *TcpClient
	bindClients    util.ConcMap
	unbindClients  util.ConcMap
	WsBroadcast    func(string, int, string, []byte)
	WsCloseHandler func(string, int, string)
	WsListenerRegister func(string, string, func(...interface {
	}), ...int)
	WsListenerUnregister func(string, string)
}

var TcpManagerInst = TcpManagerInit()
//var TcpManagerInst = TcpManagerInit()

func TcpManagerInit() *TcpManager {
	m := &TcpManager{
		bind:          make(chan *TcpClient, 10),
		unbind:        make(chan *TcpClient, 10),
		close:         make(chan *TcpClient, 10),
		bindClients:   util.NewConcMap(),
		unbindClients: util.NewConcMap(),
	}
	ants.Submit(m.Run)
	return m
}
func (tcpManager *TcpManager) GetTcpClient(appId string) (*TcpClient, bool) {
	if tmp, ok := tcpManager.bindClients.Get(appId); ok {
		tmpKeys := tmp.(util.ConcMap).Keys()
		tmpKey := util.AOrB(func() bool {
			return len(tmpKeys) > 1
		}, tmpKeys[int(time.Now().Unix())%len(tmpKeys)], tmpKeys[0]).(string)
		a, okB := tmp.(util.ConcMap).Get(tmpKey)
		if okB {
			return a.(*TcpClient), true
		} else {
			return nil, false
		}

	} else {
		return nil, false
	}
}

func (tcpManager *TcpManager) SetBroadcast(f func(string, int, string, []byte)) {
	tcpManager.WsBroadcast = f
}
func (tcpManager *TcpManager) SetCloseWsHandler(f func(string, int, string)) {
	tcpManager.WsCloseHandler = f
}
func (tcpManager *TcpManager) SetWsListenerRegister(f func(string, string, func(...interface{}), ...int)) {
	tcpManager.WsListenerRegister = f
}
func (tcpManager *TcpManager) Run() {
	for {
		select {
		case tcpClient := <-tcpManager.bind:
			var tmpM util.ConcMap
			if a, ok := tcpManager.bindClients.Get(tcpClient.appId); ok {
				tmpM = a.(util.ConcMap)
			} else {
				tmpM = util.NewConcMap()
			}
			tmpM.Set(tcpClient.ip, tcpClient)
			tcpManager.bindClients.Set(tcpClient.appId, tmpM)
		case tcpClient := <-tcpManager.unbind:
			tcpManager.unbindClients.Set(tcpClient.ip, tcpClient)
		case tcpClient := <-tcpManager.close:
			tcpManager.unbindClients.Remove(tcpClient.ip)
			if a, ok := tcpManager.bindClients.Get(tcpClient.appId); ok {
				a.(util.ConcMap).Remove(tcpClient.ip)
			}
		}
	}
}
func (tcpManager *TcpManager) SetWsMethod(WsBroadcast func(appId string, taskId int, uid string, msg []byte), CloseClient func(appId string, taskId int, uid string), Register func(appId string, ip string, f func(a ...interface{}), events ...int), Unregister func(appId string, ip string)) {
	tcpManager.WsBroadcast = WsBroadcast
	tcpManager.WsCloseHandler = CloseClient
	tcpManager.WsListenerRegister = Register
	tcpManager.WsListenerUnregister = Unregister
}
