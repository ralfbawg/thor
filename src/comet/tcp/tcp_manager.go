package tcp

import (
	"github.com/panjf2000/ants"
)

type TcpManager struct {
	bind           chan *TcpClient
	unbind         chan *TcpClient
	WsBroadcast    func(string, int, string, []byte)
	WsCloseHandler func(string, int, string)
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

func (tcpManger *TcpManager) SetBroadcast(f func(string, int, string, []byte)) {
	tcpManger.WsBroadcast = f
}
func (tcpManger *TcpManager) SetCloseWsHandler(f func(string, int, string)) {
	tcpManger.WsCloseHandler = f
}

func (tcpManger *TcpManager) Run() {
	for {
		select {
		case tcpClient := <-tcpManger.bind:
			bindClients.Set(tcpClient.appId, tcpClient)
			unbindClients.Remove(tcpClient.appId)
		case tcpClient := <-tcpManger.unbind:
			unbindClients.Set(tcpClient.appId, tcpClient)
		}
	}
}
