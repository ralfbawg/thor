package tcp

import "filter"

type TcpApp struct {
	Tasks     []*TcpTask
	wsManager *WsManager
	// App id
	appId       string
	filter      filter.WsFilterChain
	clientCount int64
	countC      chan int64
}