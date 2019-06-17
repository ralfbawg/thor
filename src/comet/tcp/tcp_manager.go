package tcp

import (
	"github.com/panjf2000/ants"
)

type TcpManager struct {
	bind                 chan *TcpClient
	unbind               chan *TcpClient
	WsBroadcast          func(string, int, string, []byte)
	WsCloseHandler       func(string, int, string)
	WsListenerRegister   func(string, string, func(...interface{}), ...int)
	WsListenerUnregister func(string, string)
}

var TcpManagerInst = TcpManagerInit()
//var TcpManagerInst = TcpManagerInit()

func TcpManagerInit() *TcpManager {
	m := &TcpManager{
		bind:   make(chan *TcpClient, 10),
		unbind: make(chan *TcpClient, 10),
	}
	ants.Submit(m.Run)
	return m
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
			bindClients.Set(tcpClient.appId, tcpClient)
			unbindClients.Remove(tcpClient.appId)
		case tcpClient := <-tcpManager.unbind:
			unbindClients.Set(tcpClient.appId, tcpClient)
		}
	}
}
func (tcpManager *TcpManager) SetWsMethod(WsBroadcast func(appId string, taskId int, uid string, msg []byte), CloseClient func(appId string, taskId int, uid string), Register func(appId string, ip string, f func(a ...interface{}), events ...int), Unregister func(appId string, ip string)) {
	tcpManager.WsBroadcast = WsBroadcast
	tcpManager.WsCloseHandler = CloseClient
	tcpManager.WsListenerRegister = Register
	tcpManager.WsListenerUnregister = Unregister
}
