package websocket

import (
	"filter"
	"common/logging"
	"time"
)

type WsApp struct {
	Tasks     []*WsTask
	wsManager *WsManager
	// App id
	appId       string
	filter      filter.WsFilterChanin
	clientCount int64
	countC      chan int64
}

func (app *WsApp) AddClient(appId string, taskId int, uid string, client *WsTaskClient) {
	task := app.Tasks[taskId]
	task.AddClient(uid, client.conn)
}

func (app *WsApp) Broadcast(msg []byte) {
	start := time.Now()
	logging.Debug("app id(%s) broadcast start %f", app.appId, start)
	for _, task := range app.Tasks {
		task.Broadcast(msg)
	}
	logging.Debug("app id(%s) broadcast end %f,cost time:%f", app.appId, time.Now(), time.Now().Sub(start).Seconds())
}

func (app *WsApp) InitFilter() bool {
	app.filter = filter.WsFilters

}
func (app *WsApp) Init() bool {
	return app.InitFilter()
}
func NewWsApp(wsManager *WsManager, appId string) (*WsApp, error) {
	app := &WsApp{
		Tasks:       make([]*WsTask, 10),
		wsManager:   wsManager,
		appId:       appId,
		clientCount: 0,
		countC:      make(chan int64, 20),
	}
	app.Init()
	return app, nil
}
