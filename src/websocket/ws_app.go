package websocket

import "util"

type WsApp struct {
	Tasks     util.ConcMap
	wsManager *WsManager
	// App id
	appId string
}

func (app *WsApp) AddClient(appId string, taskId int, client *WsTaskClient) {
	task := app.Tasks.Get(appId, taskId).(*WsTask)
	task.AddClient(appId, client.conn)
}

func (app *WsApp) Broadcast([]byte msg) {
	for task := range app.Tasks.Items() {
		task.(*WsTask).Broadcast(msg)
	}
}
