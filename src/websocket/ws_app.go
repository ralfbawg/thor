package websocket

import (
	"util"
	"filter"
)

type WsApp struct {
	Tasks     []*WsTask
	wsManager *WsManager
	// App id
	appId  string
	filter filter.WsFilterChanin
}

func (app *WsApp) AddClient(appId string, taskId int, client *WsTaskClient) {
	task := app.Tasks[taskId]
	task.AddClient(appId, client.conn)
}

func (app *WsApp) Broadcast(msg []byte) {
	for task := range app.Tasks {
		task.Broadcast(msg)
	}
}

func (app *WsApp) InitFilter() bool {
	return true
}
