package websocket

import (
	"filter"
	"common/logging"
	"time"
	"github.com/gorilla/websocket"
	rpc "comet/service/rpc"
)

type WsApp struct {
	Tasks     []*WsTask
	wsManager *WsManager
	// App id
	appId       string
	filter      filter.WsFilterChain
	clientCount int64
	countC      chan int64
}

func (app *WsApp) AddClient(taskId int, uid string, con *websocket.Conn) {
	if app.Tasks[taskId] == nil {
		app.Tasks[taskId] = NewWsTask(app, taskId)
	}
	task := app.Tasks[taskId]
	task.AddClient(uid, con)
}

func (app *WsApp) Broadcast(msg []byte) {
	start := time.Now()
	logging.Debug("app id(%s) broadcast start %f", app.appId, start)
	for _, task := range app.Tasks {
		if task != nil {
			task.Broadcast(msg)
		}

	}
	logging.Debug("app id(%s) broadcast end %f,cost time:%f", app.appId, time.Now(), time.Now().Sub(start).Seconds())
}

func (app *WsApp) InitFilter(appId string) bool {
	app.filter = filter.NewWsFilterChain(appId)
	return true
}
func (app *WsApp) Init() bool {
	return app.InitFilter(app.appId)
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
func (app *WsApp) GetAppId() string {
	return app.appId
}

func (app *WsApp) processMsg(taskId int, uid string, msg []byte) {
	if err := rpc.RpcClient.SendMsg(app.appId, taskId, uid, msg); err != nil {
		app.Tasks[taskId].GetClient(uid).Send([]byte(err.Error()))
	}
}
