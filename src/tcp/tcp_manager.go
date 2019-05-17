package tcp

type TcpManager struct {
	bind   chan TcpClient
	unbind chan TcpClient
}

var TcpManagerInst = TcpManagerInit()

func TcpManagerInit() *TcpManager {
	return &TcpManager{
		bind:   make(chan TcpClient, 10),
		unbind: make(chan TcpClient, 10),
	}
}

func (tcpManger *TcpManager) Run() {
	for {
		select {
		case tcpClient := <-tcpManger.bind:
			bindClients.Set(tcpClient.appId, tcpClient)
		case tcpClient := <-tcpManger.unbind:
			unbindClients.Set(tcpClient.appId, tcpClient)
		}
	}
}
